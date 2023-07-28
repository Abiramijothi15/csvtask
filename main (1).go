package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Record struct {
	Name  string
	Email string
	Phone string
}

func main() {
	// Create a CSV file
	file, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write some data to the CSV file
	writer := csv.NewWriter(file)
	writer.Write([]string{"Name", "Email", "Phone"})
	writer.Write([]string{"Abirami", "abi.doe@example.com", "123-456-7890"})
	writer.Write([]string{"Tiger", "tig.smith@example.com", "555-555-5555"})
	writer.Flush()

	// Read the data from the CSV file
	csvFile, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Convert the data to JSON format
	var jsonData []Record
	for _, record := range records[1:] {
		jsonData = append(jsonData, Record{
			Name:  record[0],
			Email: record[1],
			Phone: record[2],
		})
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	// Write the JSON data to a file
	err = ioutil.WriteFile("data.json", jsonBytes, 0644)
	if err != nil {
		panic(err)
	}

	// Insert the data into a MySQL database
	db, err := sql.Open("mysql", "root:mypassword@tcp(69.195.138.156)/test")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO people (name, email, phone) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}

	for _, record := range jsonData {
		_, err = stmt.Exec(record.Name, record.Email, record.Phone)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Data successfully inserted into MySQL database")
}
