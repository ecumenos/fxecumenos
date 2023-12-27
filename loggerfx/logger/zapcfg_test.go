package logger_test

import (
	"log"

	"github.com/ecumenos/fxecumenos/loggerfx/logger"
)

func ExampleNewDevelopmentLogger() {
	serviceName := "MyService"
	devLogger, err := logger.NewDevelopmentLogger(serviceName)
	if err != nil {
		log.Fatal(err)
	}
	sugaredLogger := devLogger.Sugar()

	sugaredLogger.Infow("Log the logs", "logName", "loggy", "logValue", 1.2)

	// Outputs to stderr: {"level":"info","time":"2018-12-11T21:17:06.128Z","caller":"zapconf/zapconf_test.go:28","msg":"Log the logs","service":"MyService","logName":"loggy","logValue":1.2}
}

func ExampleNewProductionLogger() {
	serviceName := "MyService"
	prodLogger, err := logger.NewProductionLogger(serviceName)
	if err != nil {
		log.Fatal(err)
	}
	sugaredLogger := prodLogger.Sugar()

	sugaredLogger.Errorw("Log the logs", "logName", "loggy", "logValue", 1.2)

	// Outputs to stderr: {"level":"error","time":"2018-12-11T21:17:06.128Z","caller":"zapconf/zapconf_test.go:28","msg":"Log the logs","service":"MyService","logName":"loggy","logValue":1.2}
}
