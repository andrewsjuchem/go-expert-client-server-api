package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CurrencyExchangeQuote struct {
	Bid string `json:"bid"`
}

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

func main() {
	quote, err := GetCurrencyExchange()
	if err != nil {
		Sugar.Error(err)
		panic(err)
	}
	saveQuoteToFile(quote)
}

func GetCurrencyExchange() (*CurrencyExchangeQuote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Prepare request
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		Sugar.Error(err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Run request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Sugar.Error(err)
		return nil, err
	}

	// Parse response
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		Sugar.Error(err)
		return nil, err
	}

	var quote CurrencyExchangeQuote
	// Convert json response into struct object
	err = json.Unmarshal(body, &quote)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func saveQuoteToFile(quote *CurrencyExchangeQuote) error {
	// Create the log folder if it does not exist
	filePath := "./../outputs/"
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		Sugar.Error(err)
		return err
	}

	// Write the quote to the file
	err = os.WriteFile(filePath+"cotacao.txt", []byte("Dólar: "+quote.Bid), 0644)
	if err != nil {
		Sugar.Error(err)
		return err
	}
	return nil
}
