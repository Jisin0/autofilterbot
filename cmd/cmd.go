package cmd

import (
	"flag"

	"github.com/Jisin0/autofilterbot/internal/core"
)

// Execute acts as the entry point of the application, it parses command line arguments and then runs the application.
func Execute() {
	mongodbUri := flag.String("mongodb-uri", "", "mongodb uri for database (use env for additional uris)")
	botToken := flag.String("bot-token", "", "bot token obtained from @botfather")
	logLevel := flag.String("log-level", "info", "level of logs to be shown")
	noOutput := flag.Bool("no-output", false, "disable console logs (does not affect log file)")

	flag.Parse()

	core.Run(
		core.RunAppOptions{
			MongodbURI:         *mongodbUri,
			LogLevel:           *logLevel,
			BotToken:           *botToken,
			DisableConsoleLogs: *noOutput,
		},
	)
}
