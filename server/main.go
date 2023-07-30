package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/mattn/go-sqlite3"
)

type CurrencyExchangeQuote struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type CurrencyExchangeQuotePublic struct {
	Bid string `json:"bid"`
}

func (q *CurrencyExchangeQuote) ToPublic() *CurrencyExchangeQuotePublic {
	return &CurrencyExchangeQuotePublic{
		Bid: q.USDBRL.Bid,
	}
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

func init() {
	InitializeLogger()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", CurrencyExchangeHandler)
	Sugar.Info("Listing to port 8080")
	http.ListenAndServe(":8080", mux)
}

func CurrencyExchangeHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := GetCurrencyExchange()
	quotePublic := quote.ToPublic()
	if err != nil {
		Sugar.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotePublic)
}

func GetCurrencyExchange() (*CurrencyExchangeQuote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Prepare request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

	err = insertQuote(&quote)
	if err != nil {
		Sugar.Error(err)
		return nil, err
	}

	return &quote, nil
}

func insertQuote(quote *CurrencyExchangeQuote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Specify the path and name of the database file
	dbPath := "./../databases/sqlite3"

	// Create the folder if it does not exist
	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		Sugar.Error(err)
		return err
	}

	// Open the database file
	db, err := sql.Open("sqlite3", dbPath+"/currency_exchange.db")
	if err != nil {
		Sugar.Error(err)
		return err
	}
	defer db.Close()

	// Create table if not exist
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS quote(id INTEGER PRIMARY KEY, code TEXT NOT NULL, codein TEXT NOT NULL, exchange_rate NUMERIC NOT NULL, create_date TEXT NOT NULL)", nil)
	if err != nil {
		Sugar.Error(err)
		return err
	}

	// Insert a row into the table
	_, err = db.ExecContext(ctx, "INSERT INTO quote (code, codein, exchange_rate, create_date) VALUES (?, ?, ?, ?)",
		quote.USDBRL.Code,
		quote.USDBRL.Codein,
		quote.USDBRL.Bid,
		quote.USDBRL.CreateDate)
	if err != nil {
		Sugar.Error(err)
		return err
	}
	return nil
}
