package entity

import (
	"fmt"
	"math/rand"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func CreateMockBatch() []User {
	var users []User

	for i := 1; i <= 1000; i++ {
		userName := fmt.Sprintf("User_%d", i)
		users = append(users, User{Name: userName, Age: rand.Intn(34)})
	}

	return users
}
