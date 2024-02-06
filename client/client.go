package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const requestExpirationTime = 300 * time.Millisecond

type Cotation struct {
	From      string `json:"code"`
	To        string `json:"codein"`
	Name      string `json:"name"`
	High      string `json:"high"`
	Low       string `json:"low"`
	VarBid    string `json:"varBid"`
	PctChange string `json:"pctChange"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	CreatedAt string `json:"-"`
}

type errorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func main() {

	cotation, err := GetCotation("usd", "brl")
	writeTofile("cotacoes", cotation)
	if err != nil {
		panic(err)
	}

}

func GetCotation(from string, to string) (string, error) {

	ctx := context.Background()
	// o contexto expira em 1 segundo!
	ctx, cancel := context.WithTimeout(ctx, requestExpirationTime)
	defer cancel() // de alguma forma nosso contexto será cancelado
	url := fmt.Sprintf("http://localhost:8080/cotacao?from=%s&to=%s", from, to)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		var businessError errorResponse
		err = json.Unmarshal(body, &businessError)
		if err != nil {
			return "", err
		}
		return "", errors.New(businessError.Message)
	} else {

		return string(body), nil

	}

}

func writeTofile(fileName string, message string) {
	// Open the file in append mode
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	f, err := os.OpenFile(fileName+".txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		panic(err)
	}
	defer f.Close()

	// Skip one line
	_, err = f.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing newline:", err)
		panic(err)
	}

	// Write the new data
	_, err = f.Write([]byte(timestamp + " - Dolar: " + message))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		panic(err)
	}
}
