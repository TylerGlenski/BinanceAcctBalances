package main

import (
        "fmt"
		"context"
		"os"
		"strconv"
		"encoding/csv"
		"log"
		"github.com/adshao/go-binance"
)


type Listing struct{
	// holds account crypto listing data
	Symbol string
	Locked string
	Free string
	Total string
}
func getAccountInfo(apiKey string, secretKey string) ([]Listing, error){
	// returns a map data structure with assets mapping to three values free locked total
	// returns a slice of listings
	
	returnValue := []Listing{}

	client := binance.NewClient(apiKey, secretKey)
	client.NewSetServerTimeService().Do(context.Background())
	res, err := client.NewGetAccountService().Do(context.Background())

	checkError("Binance api call error", err)

	for _, v := range res.Balances {
		lockedBalance, err := strconv.ParseFloat(v.Locked, 32)
		freeBalance, err := strconv.ParseFloat(v.Free, 32)
		
		checkError("locked + free var creation error: ", err)

		if lockedBalance > 0 || freeBalance > 0 {
			var total float64 = float64(lockedBalance + freeBalance)
			totalStr := strconv.FormatFloat(total, 'f', 8, 64)
			returnValue = append(returnValue, Listing{
				Symbol: v.Asset,
				Locked: v.Locked,
				Free: v.Free,
				Total: totalStr,
			})
		}

	}

	return returnValue, nil
}


func checkError(message string, err error) {
	// Checks for an error, raises error if error
	// no return value
    if err != nil {
        log.Fatal(message, err)
    }
}


func writeToCSV(balances []Listing) {
	// Takes in a slice of listings (accountBalances in main)
	// Writes data to a csv file called result.csv
	// no return value

	var data = [][]string{{"Symbol", "Free", "Locked", "Total"}}

	for _, v := range balances{
		appendLine := []string {v.Symbol, v.Free, v.Locked, v.Total}
		data = append(data, appendLine)
	}

    file, err := os.Create("result.csv")
    checkError("Cannot create file", err)
	defer file.Close()
	
    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, value := range data {
        err := writer.Write(value)
        checkError("Cannot write to file", err)
	}
	fmt.Println("Writing successfull located at ./result.csv")
}


func main() {
    // Makes an api call to binance api and writes result to a csv
	
	// PUT YOUR API KEY/SECRET HERE! A better way would be to read from env, or a secret on AWS.
	apiKey := ""
	apiSecret := ""
	
	accountBalances, err := getAccountInfo(apiKey, apiSecret)
	checkError("CRITICAL ERROR: ", err)
	writeToCSV(accountBalances)

}