package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/andrewsjuchem/go-expert-client-server-api/utils"
)

type CurrencyExchangeQuote struct {
	Bid string `json:"bid"`
}

func init() {
	utils.InitializeLogger()
}

func main() {
	utils.Sugar.Info("Getting currency exchange quote")
	quote, err := GetCurrencyExchange()
	if err != nil {
		utils.Sugar.Error(err)
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

	return &quote, nil
}

func saveQuoteToFile(quote *CurrencyExchangeQuote) error {
	// Create the log folder if it does not exist
	filePath := "./../outputs/"
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		utils.Sugar.Error(err)
		return err
	}

	// Write the quote to the file
	err = os.WriteFile(filePath+"cotacao.txt", []byte("Dólar: "+quote.Bid), 0644)
	if err != nil {
		utils.Sugar.Error(err)
		return err
	}
	return nil
}
