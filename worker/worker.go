package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
)

const API_URL = "http://localhost:3000/api/process"

type Request struct {
	UserID int `json:"userId"`
}

type Response struct {
	UserID int `json:"userId"`
}

func processUserIDs(userIDs []int, wg *sync.WaitGroup, results chan<- Response) {
	defer wg.Done()
	for _, userID := range userIDs {
		reqBody, err := json.Marshal(Request{UserID: userID})
		if err != nil {
			fmt.Println("Error marshalling request:", err)
			continue
		}

		resp, err := http.Post(API_URL, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		var response Response
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error unmarshalling response:", err)
			continue
		}

		results <- response
	}
}

func main() {
	file, err := os.Open("../user_ids.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	var userIDs []int
	csvReader := csv.NewReader(file)
	_, _ = csvReader.Read() 
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		userID, _ := strconv.Atoi(record[0])
		userIDs = append(userIDs, userID)
	}

	batchSize := len(userIDs) / 10
	var wg sync.WaitGroup
	results := make(chan Response, len(userIDs))

	for i := 0; i < 10; i++ {
		start := i * batchSize
		end := start + batchSize
		if i == 9 {
			end = len(userIDs) 
		}
		wg.Add(1)
		go processUserIDs(userIDs[start:end], &wg, results)
	}

	wg.Wait()
	close(results)

	processedUserIDs := make([]int, 0, len(userIDs))
	for result := range results {
		processedUserIDs = append(processedUserIDs, result.UserID)
	}

	sort.Ints(processedUserIDs)

	outputFile, err := os.Create("../processed_user_ids.csv")
	if err != nil {
		fmt.Println("Error creating output CSV file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	writer.Write([]string{"userId"}) 

	for _, userID := range processedUserIDs {
		writer.Write([]string{strconv.Itoa(userID)})
	}

	fmt.Println("New CSV file has been created with processed user IDs")
}
