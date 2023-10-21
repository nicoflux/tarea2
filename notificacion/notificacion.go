package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	client := &http.Client{}
	data, err := os.Open("order.json")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	req, err := http.NewRequest("POST", "https://sjwc0tz9e4.execute-api.us-east-2.amazonaws.com/Prod", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}
