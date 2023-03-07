package main

import (
	"fmt"
	"math/rand"
	"time"
)

var db *Storage

type Person struct {
	FirstName string
	LastName  string
	Age       int
	City      string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var lastRecord Person

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateData(amount int) (persons []Person) {
	persons = make([]Person, amount)
	for i := 0; i < amount; i++ {
		persons[i] = Person{randomString(4), randomString(3), rand.Int(), "Ottawa-" + randomString(2)}
	}
	return persons
}

func writeData(records *[]Person) {
	primaries := []string{"FirstName", "LastName"}
	var err error

	for _, x := range *records {
		err = db.Write(GetIdentifier(x, primaries), x)
		if err != nil {
			panic(err)
		}
	}

	lastRecord = (*records)[len(*records)-1]
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialization
	db = &Storage{}
	db.Init()

	// Single write example
	start := time.Now()
	x := Person{"Daniil", "Furmanov", 10, "Moscow"}
	db.Write(GetIdentifier(x, nil), x)
	elapsed := time.Since(start)

	fmt.Printf("Single write took %s\n", elapsed)

	// Generate large amount of data
	fmt.Println("Generating data...")
	start = time.Now()
	data := generateData(1000000)
	elapsed = time.Since(start)

	fmt.Printf("Generation took %s\n", elapsed)

	// Write this data
	fmt.Println("Writing the data...")
	start = time.Now()
	writeData(&data)
	elapsed = time.Since(start)

	fmt.Printf("Database now contains %d records. Writing took %s\n", db.Size(), elapsed)

	// Search Example
	searchPattern := map[string]interface{}{
		"FirstName": lastRecord.FirstName,
		"LastName":  lastRecord.LastName,
	}

	segment := GetIdentifierFromMap(searchPattern, []string{"FirstName", "LastName"})

	start = time.Now()

	fmt.Println("Searching data... segment " + segment)

	seg, err := db.FindSegment(segment)
	if err != nil {
		fmt.Println("Segment not found")
		return
	}

	fmt.Println("Segment found. Path:", seg.Path, "; Size:", seg.Size, "; Records count:", len(seg.Records))

	records := seg.FindRecords(searchPattern)
	elapsed = time.Since(start)

	if len(records) > 0 {
		fmt.Println(records[0].Data)
	} else {
		fmt.Println("Zero results found")
	}

	fmt.Printf("Search took %s\n", elapsed)

}
