package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Jisin0/autofilterbot/internal/config"
	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/database/mongo"
	"github.com/Jisin0/autofilterbot/pkg/autodelete"
	"github.com/Jisin0/autofilterbot/pkg/log"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var _app *App

// App wraps various individual components of the app to orchestrate application processes.
type App struct {
	DB        database.Database
	Log       *zap.Logger
	StartTime time.Time

	Config     config.Config
	AutoDelete *autodelete.Manager
}

// RunAppOptions wraps command-line arguments for app startup.
type RunAppOptions struct {
	MongodbURI string
	LogLevel   string
	BotToken   string
}

// Run starts the application and initializes core components.
func Run(opts RunAppOptions) {
	err := godotenv.Load("config.env", ".env") // config.env is supported bcuz other repos use it for some reason
	if err != nil {
		fmt.Println("ERROR: load variables from .env file failed", err)
	}

	logLevel := opts.LogLevel
	if s := os.Getenv("LOG_LEVEL"); s != "" {
		logLevel = s
	}

	log.Initialize(logLevel)
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
		Config:     *appConfig,
		Log:        logger,
		AutoDelete: autodeleteManager,
		StartTime:  time.Now(),
	}

	//TODO: setup dispatcher and run bot

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c // wait until an interrupt signal is received
	logger.Info("stopping app: interrupt signal received", zap.Any("signal", s))

	cancel() // autodelete & mongo updater should stop with this
	_app.DB.Shutdown()
}

// App returns the initialized global app instance.
func Application() *App {
	return _app
}
