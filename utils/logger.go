package utils

import (
	"fmt"
	"log"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func InitializeLogger(processName string) {
	// Create a logger configuration
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	defaultLogLevel := zapcore.DebugLevel

	// Console encoder so it prints the logs to the console
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Create the log folder if it does not exist
	workingDirectory, _ := os.Getwd()
	logDirectory := path.Join(workingDirectory, "./logs")
	// if os.Getenv("APP_ENV") == "docker" {
	// 	logDirectory = path.Join(workingDirectory, "./logs")
	// } else {
	// 	logDirectory = path.Join(workingDirectory, "./../logs")
	// }
	var logFileName string
	if len(processName) > 0 {
		logFileName = fmt.Sprintf(logDirectory+"/log_%s_%d.log", processName, os.Getpid())
	} else {
		logFileName = fmt.Sprintf(logDirectory+"/log_%d.log", os.Getpid())
	}
	err := os.MkdirAll(logDirectory, os.ModePerm)
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
