package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Create("user_ids.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"userId"}) // Write header

	for i := 1; i <= 10000; i++ {
		writer.Write([]string{strconv.Itoa(i)})
	}

	fmt.Println("CSV file has been created with 10000 user IDs")
}
