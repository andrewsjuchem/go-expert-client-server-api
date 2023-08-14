package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/andrewsjuchem/go-expert-client-server-api/utils"
)

type CurrencyExchangeQuote struct {
	Bid string `json:"bid"`
}

func init() {
	utils.InitializeLogger("client")
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
	var endpointUrl string
	if os.Getenv("APP_ENV") == "docker" {
		endpointUrl = "http://go-server-api:8080/cotacao"
	} else {
		endpointUrl = "http://localhost:8080/cotacao"
	}
	req, err := http.NewRequestWithContext(ctx, "GET", endpointUrl, nil)
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
	var fileDirectory string
	workingDirectory, _ := os.Getwd()
	if os.Getenv("APP_ENV") == "docker" {
		fileDirectory = path.Join(workingDirectory, "./outputs")
	} else {
		fileDirectory = path.Join(workingDirectory, "./../outputs/")
	}
	err := os.MkdirAll(fileDirectory, os.ModePerm)
	if err != nil {
		utils.Sugar.Error(err)
		return err
	}

	// Write the quote to the file
	err = os.WriteFile(fileDirectory+"/cotacao.txt", []byte("Dólar: "+quote.Bid), 0644)
	if err != nil {
		utils.Sugar.Error(err)
		return err
	}
	return nil
}
