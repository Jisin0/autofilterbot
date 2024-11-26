package cmd

import (
	"flag"

	"github.com/Jisin0/autofilterbot/internal/app"
)

// Execute acts as the entry point of the application, it parses command line arguments and then runs the application.
func Execute() {
	mongodbUri := flag.String("mongodb-uri", "", "mongodb uri for database (use env for additional uris or couchbase)")
	botToken := flag.String("bot-token", "", "bot token obtained from @botfather")
	logLevel := flag.String("log-level", "warn", "level of logs to be shown")

	app.Run(
		app.RunAppOptions{
			MongodbURI: *mongodbUri,
			LogLevel:   *logLevel,
			BotToken:   *botToken,
		},
	)
}
