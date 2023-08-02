package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/andrewsjuchem/go-expert-client-server-api/utils"
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

func init() {
	utils.InitializeLogger()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", CurrencyExchangeHandler)
	utils.Sugar.Info("Listing to port 8080")
	http.ListenAndServe(":8080", mux)
}

func CurrencyExchangeHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := GetCurrencyExchange()
	quotePublic := quote.ToPublic()
	if err != nil {
		utils.Sugar.Error(err)
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
		utils.Sugar.Error(err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Run request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Sugar.Error(err)
		return nil, err
	}

	// Parse response
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		utils.Sugar.Error(err)
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
		utils.Sugar.Error(err)
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
		utils.Sugar.Error(err)
		return err
	}

	// Open the database file
	db, err := sql.Open("sqlite3", dbPath+"/currency_exchange.db")
	if err != nil {
		utils.Sugar.Error(err)
		return err
	}
	defer db.Close()

	// Create table if not exist
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS quote(id INTEGER PRIMARY KEY, code TEXT NOT NULL, codein TEXT NOT NULL, exchange_rate NUMERIC NOT NULL, create_date TEXT NOT NULL)", nil)
	if err != nil {
		utils.Sugar.Error(err)
		return err
	}

	// Insert a row into the table
	_, err = db.ExecContext(ctx, "INSERT INTO quote (code, codein, exchange_rate, create_date) VALUES (?, ?, ?, ?)",
		quote.USDBRL.Code,
		quote.USDBRL.Codein,
		quote.USDBRL.Bid,
		quote.USDBRL.CreateDate)
	if err != nil {
		utils.Sugar.Error(err)
		return err
	}
	return nil
}
