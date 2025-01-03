package core

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Jisin0/autofilterbot/internal/config"
	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/database/mongo"
	"github.com/Jisin0/autofilterbot/pkg/autodelete"
	exthandlers "github.com/Jisin0/autofilterbot/pkg/filters"
	"github.com/Jisin0/autofilterbot/pkg/log"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var _app *App

// App wraps various individual components of the app to orchestrate application processes.
type App struct {
	DB        database.Database
	Log       *zap.Logger
	StartTime time.Time
	Bot       *gotgbot.Bot

	Config     *config.Config
	AutoDelete *autodelete.Manager
}

// extendedHandler returns a handlers.Response that calls

// RunAppOptions wraps command-line arguments for app startup.
type RunAppOptions struct {
	MongodbURI         string
	LogLevel           string
	BotToken           string
	DisableConsoleLogs bool
}

// Run starts the application and initializes core components.
func Run(opts RunAppOptions) {
	err := godotenv.Load(".env") // config.env is supported bcuz other repos use it for some reason
	if err != nil {
		fmt.Println("ERROR: load variables from .env file failed", err)
	}

	logLevel := opts.LogLevel
	if s := os.Getenv("LOG_LEVEL"); s != "" {
		logLevel = s
	}

	log.Initialize(logLevel, opts.DisableConsoleLogs)
	logger := log.Logger()

	botToken := opts.BotToken
	if s := os.Getenv("BOT_TOKEN"); s != "" {
		botToken = s
	}

	if botToken == "" {
		logger.Fatal("bot token not provided")
	}

	bot, err := gotgbot.NewBot(botToken, &gotgbot.BotOpts{})
	if err != nil {
		logger.Fatal("create bot failed", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background()) // all background jobs (tickers) must use this ctx

	mongodbUri := opts.MongodbURI
	if s := os.Getenv("MONGODB_URI"); s != "" {
		mongodbUri = s
	}

	couchBaseUri := os.Getenv("COUCHBASE_URI")
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	var db database.Database

	switch {
	case mongodbUri != "":
		var additionalUri []string
		for i := 1; i <= 5; i++ { // attempts to fetch MONGODB_URI1 to MONGODB_URI5. //TODO: remove hardcoded limit after testing
			if s := os.Getenv(fmt.Sprintf("MONGODB_URI%d", i)); s != "" {
				additionalUri = append(additionalUri, s)
			}
		}

		db, err = mongo.NewClient(ctx, mongodbUri, databaseName, collectionName, additionalUri...)
	case couchBaseUri != "":
		logger.Fatal("not implemented")
	default:
		logger.Fatal("mongodb or couchbase uri not found, please read the database setup guide")
	}

	if err != nil {
		logger.Fatal("database setup failed", zap.Error(err))
	}

	appConfig, err := db.GetConfig(bot.Id)
	if err != nil {
		logger.Error("failed to load configs from db", zap.Error(err))
	}

	autodeleteManager, err := autodelete.NewManager(bot)
	if err != nil {
		logger.Error("autodelete module setup failed", zap.Error(err))
	}

	_app = &App{
		DB:         db,
		Config:     appConfig,
		Bot:        bot,
		Log:        logger,
		AutoDelete: autodeleteManager,
		StartTime:  time.Now(),
	}

	dispatcher := SetupDispatcher(logger)
	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{})

	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
	})
	if err != nil {
		logger.Fatal(
			"failed to start polling updates",
			zap.Error(err),
		)
	}

	logger.Info(fmt.Sprintf("@%s started successfully !", bot.Username))

	//TODO: setup dispatcher and run bot
	updater.Idle()
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)

	// s := <-c // wait until an interrupt signal is received
	logger.Info("stopping app: interrupt signal received", zap.Any("signal", ""))

	updater.Stop()

	cancel() // autodelete & mongo updater should stop with this
	_app.DB.Shutdown()
}

const (
	commandHandlerGroup = 1
)

// SetupDispatcher creates a new empty dispatcher with error and panic recovery setup.
func SetupDispatcher(log *zap.Logger) *ext.Dispatcher {
	d := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logFields := []zap.Field{zap.Error(err)}
			logFields = addLogFieldsFromContext(ctx, logFields)

			log.Error("error occurred while handling update", logFields...)

			return ext.DispatcherActionNoop
		},

		Panic: func(b *gotgbot.Bot, ctx *ext.Context, r interface{}) {
			logFields := []zap.Field{zap.String("panic", fmt.Sprintf("%v\n%s", r, cleanedStack()))}
			logFields = addLogFieldsFromContext(ctx, logFields)

			log.Error("panic recovered", logFields...)
		},
	})

	d.AddHandlerToGroup(exthandlers.NewCommands([]string{"start", "about", "help", "privacy"}, StaticCommands), commandHandlerGroup)

	return d
}

// logFieldsFromContext adds zap fields to logFields about specific update from ctx to help troubleshooting.
func addLogFieldsFromContext(ctx *ext.Context, logFields []zap.Field) []zap.Field {
	switch {
	case ctx.Message != nil:
		logFields = append(logFields,
			zap.Int64("chat_id", ctx.Message.Chat.Id),
			zap.Int64("message_id", ctx.Message.MessageId),
			zap.String("text", ctx.Message.Text),
		)
	case ctx.CallbackQuery != nil:
		logFields = append(logFields,
			zap.String("callback_query_id", ctx.CallbackQuery.Id),
			zap.String("data", ctx.CallbackQuery.Data),
			zap.Int64("chat_id", ctx.CallbackQuery.Message.GetChat().Id),
			zap.Int64("message_id", ctx.CallbackQuery.Message.GetMessageId()),
		)
	case ctx.InlineQuery != nil:
		logFields = append(logFields,
			zap.String("inline_query_id", ctx.InlineQuery.Id),
			zap.String("query", ctx.InlineQuery.Query),
			zap.Int64("from", ctx.InlineQuery.From.Id),
		)
	}

	return logFields
}

// cleanedStack returns stack trace with gotgbot library parts removed to prevent confusion.
// Copied from https://github.com/PaulSonOfLars/gotgbot/blob/v2/ext/dispatcher.go.
func cleanedStack() string {
	return strings.Join(strings.Split(string(debug.Stack()), "\n")[4:], "\n")
}

// App returns the initialized global app instance.
func Application() *App {
	return _app
}
