package main

import (
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// start with a raw JSON string
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"initialFields": {"pid": 12345},
		 "encoderConfig": {
	     "messageKey": "message",
	     "levelKey": "level",
	     "levelEncoder": "lowercase"
		}
	}`)

	var cfg zap.Config
	// Unmarshal the JSON into a zap.Config struct
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.DisableCaller = false
	//custom time encoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	// Don't forget to flush the logger
	defer logger.Sync()

	additionalField1 := zap.String("url", "codewithmukesh.com")
	additionalField2 := zap.Int("attempt", 3)
	logger.Info("This is an INFO message with additional fields.", additionalField1, additionalField2)
	logger.Error("This is an ERROR message with additional fields.", additionalField1, additionalField2)
}
