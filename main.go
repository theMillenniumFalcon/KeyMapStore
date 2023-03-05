package main

import "math/rand"

var DB *Storage

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
