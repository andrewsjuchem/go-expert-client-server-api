package utils

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func InitializeLogger() {
	// Create a logger configuration
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	defaultLogLevel := zapcore.DebugLevel

	// Console encoder so it prints the logs to the console
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Create the log folder if it does not exist
	logPath := "./../logs/"
	logFileName := fmt.Sprintf(logPath+"log_%d.log", os.Getpid())
	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// File encoder so it prints the logs to a json file
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, _ := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)

	// Create a logger instance
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Sugar = Logger.Sugar()
	defer Logger.Sync()
}

// func InitializeLogger() {
// 	rawJSON := []byte(`{
// 			"level": "debug",
// 			"development": false,
// 			"disableCaller": false,
// 			"disableStacktrace": false,
// 			"encoding": "json",
// 			"outputPaths": ["stdout"],
// 			"errorOutputPaths": ["stderr"],
// 			"encoderConfig": {
// 				"levelKey": "level",
// 				"messageKey": "message",
// 				"levelEncoder": "lowercase"
// 			}
// 		}`)
// 	var config zap.Config
// 	if err := json.Unmarshal(rawJSON, &config); err != nil {
// 		log.Fatal(err)
// 	}
// 	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	Logger, err := config.Build()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	Sugar = Logger.Sugar()
// 	defer Logger.Sync()
// }
