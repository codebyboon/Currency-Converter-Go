package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ExchangeRates struct {
	Success bool               `json:"success"`
	Base    string             `json:"base"`
	Rates   map[string]float64 `json:"rates"`
}

func getExchangeRates(baseCurrency string, convertCurrency string) (*ExchangeRates, error) {

	endpoint := "latest"
	apiKey := "817ba03f35aa5fa569cec60f93f11e0f"
	url := fmt.Sprintf("http://data.fixer.io/api/%s?access_key=%s&symbols=%s", endpoint, apiKey, convertCurrency)
	fmt.Println(url)

	client := &http.Client{}
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	fmt.Println("successfully retrieved API response.")
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var exchangeRates ExchangeRates
	err = json.Unmarshal(data, &exchangeRates)
	if err != nil {
		return nil, err
	}

	if !exchangeRates.Success {
		return nil, fmt.Errorf("API request failed")
	}

	return &exchangeRates, nil
}

func convertCurrency(amount float64, fromCurrency string, toCurrency string) (float64, error) {

	rates, err := getExchangeRates(fromCurrency, toCurrency)
	if err != nil {
		return 0, err
	}

	exchangeRate, found := rates.Rates[toCurrency]
	if !found {

		return 0, fmt.Errorf("conversion rate not found for %s to %s", fromCurrency, toCurrency)
	}

	convertedAmount := amount * exchangeRate

	return convertedAmount, nil
}

func main() {
	amount := 100.0
	fromCurrency := "EUR"
	toCurrency := "USD"

	convertedAmount, err := convertCurrency(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, convertedAmount, toCurrency)
}
