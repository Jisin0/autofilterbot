package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Jisin0/autofilterbot/internal/app"
	"github.com/Jisin0/autofilterbot/internal/cache"
	"github.com/Jisin0/autofilterbot/internal/configpanel"
	"github.com/Jisin0/autofilterbot/internal/database/mongo"
	"github.com/Jisin0/autofilterbot/internal/index"
	"github.com/Jisin0/autofilterbot/pkg/autodelete"
	"github.com/Jisin0/autofilterbot/pkg/env"
	"github.com/Jisin0/autofilterbot/pkg/log"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var _app *Core

// Core wraps various individual components of the app to orchestrate application processes.
type Core struct {
	app.App
	Ctx context.Context

	additionalURLsCount int
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

	if mongodbUri == "" {
		logger.Fatal("MONGODB_URI not provided. Set it as an environment variable or provide as a command line argument.")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	var additionalUri []string

	for i := 1; i <= 5; i++ { // attempts to fetch MONGODB_URI1 to MONGODB_URI5. //TODO: remove hardcoded limit after testing
		if s := os.Getenv(fmt.Sprintf("MONGODB_URI%d", i)); s != "" {
			additionalUri = append(additionalUri, s)
		}
	}

	mongoOpts := mongo.NewClientOpts{
		DatabaseName:        databaseName,
		FilesCollectionName: collectionName,
		AdditionalURLs:      additionalUri,
	}

	db, err := mongo.NewClient(ctx, mongodbUri, logger, mongoOpts)

	if err != nil {
		logger.Fatal("database setup failed", zap.Error(err))
	}

	appConfig, err := db.GetConfig(bot.Id)
	if err != nil {
		logger.Error("failed to load configs from db", zap.Error(err))
	}

	if appConfig.FileCollectionIndex != 0 {
		err = db.UpdateStorageCollection(appConfig.FileCollectionIndex)
		if err != nil {
			logger.Warn("setting custom storage collection failed, using default database", zap.Error(err))
		}
	}

	autodeleteManager, err := autodelete.NewManager(bot)
	if err != nil {
		logger.Error("autodelete module setup failed", zap.Error(err))
	}

	go autodeleteManager.Run(ctx, logger)

	_app = &Core{
		App: app.App{
			DB:           db,
			Config:       appConfig,
			Bot:          bot,
			Log:          logger,
			AutoDelete:   autodeleteManager,
			StartTime:    time.Now(),
			Cache:        cache.NewCache(),
			Admins:       env.Int64s("ADMINS"),
			IndexManager: index.NewManager(),
		},
		Ctx: ctx,
	}

	_app.ConfigPanel = configpanel.CreatePanel(_app)

	dispatcher := SetupDispatcher(logger)
	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		UnhandledErrFunc: func(err error) {
			logger.Debug("updater: unhandled error", zap.Error(err))
		},
	})

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

	go _app.RestartActiveIndexOperations(ctx)

	if appConfig.FileCollectionUpdater {
		_app.DB.RunCollectionUpdater(ctx, logger)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c // wait until an interrupt signal is received
	logger.Info("stopping app: interrupt signal received", zap.Any("signal", s))

	updater.Stop()

	cancel() // autodelete & mongo updater should stop with this
	_app.DB.Shutdown()
}

// AuthAdmin reports whether the user who sent the message is an admin or otherwise sends a warn message.
func (core *Core) AuthAdmin(ctx *ext.Context) bool {
	switch {
	case ctx.Message != nil:
		if !containsI64(_app.Admins, ctx.Message.From.Id) {
			ctx.Message.Reply(core.Bot, "<b>𝖮𝗇𝗅𝗒 𝖺𝗇 𝖺𝖽𝗆𝗂𝗇 𝖼𝖺𝗇 𝗎𝗌𝖾 𝗍𝗁𝖺𝗍 𝖼𝗈𝗆𝗆𝖺𝗇𝖽, 𝖯𝖾𝖺𝗌𝖺𝗇𝗍❗</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
			return false
		}
	case ctx.CallbackQuery != nil:
		if !containsI64(_app.Admins, ctx.CallbackQuery.From.Id) {
			ctx.CallbackQuery.Answer(_app.Bot, &gotgbot.AnswerCallbackQueryOpts{Text: "𝖮𝗇𝗅𝗒 𝖺𝗇 𝖺𝖽𝗆𝗂𝗇 𝖼𝖺𝗇 𝗎𝗌𝖾 𝗍𝗁𝖺𝗍 𝖼𝗈𝗆𝗆𝖺𝗇𝖽, 𝖯𝖾𝖺𝗌𝖺𝗇𝗍❗", ShowAlert: true})
			return false
		}
	default:
		_app.Log.Warn("authadmin: unsupported update received", zap.Int64("update_id", ctx.UpdateId))
		return false
	}

	return true
}

// RefreshConfig refetches the bot configs from db.
func (core *Core) RefreshConfig() {
	c, err := core.DB.GetConfig(core.Bot.Id)
	if err != nil {
		core.Log.Error("failed to refresh configs", zap.Error(err))
	}

	core.Config = c
}

// RestartActiveIndexOperations restarts all active index operations.
func (c *Core) RestartActiveIndexOperations(ctx context.Context) {
	ops, err := c.DB.GetActiveIndexOperations()
	if err != nil {
		_app.Log.Debug("core: failed to fetch active index operations", zap.Error(err))
		return
	}

	if len(ops) == 0 {
		return
	}

	c.Log.Debug("core: restarting active index operations", zap.Int("num", len(ops)))

	for _, i := range ops {
		ctx, o := c.IndexManager.NewOperation(ctx, i, c.DB, c.Log, c.Bot)
		c.IndexManager.RunOperation(ctx, o)
	}
}

// GetAdditionalCollectionCount returns the number of additional db urls provided.
func (c *Core) GetAdditionalCollectionCount() int {
	return c.additionalURLsCount
}

func (c *Core) SetCollectionIndex(index int) {
	err := c.DB.UpdateStorageCollection(index)
	if err != nil {
		c.Log.Warn("core: failed to update collection index", zap.Error(err))
	}
}

// App returns the initialized global app instance.
func Application() *Core {
	return _app
}

// containsI64 reports whether the given value is in a slice.
func containsI64(s []int64, val int64) bool {
	for _, i := range s {
		if i == val {
			return true
		}
	}
	return false
}
