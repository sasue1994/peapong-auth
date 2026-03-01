package main

import (
	"fmt"
	"testing"
)

func TestHashPassword(test *testing.T) {

	password := "testPassword"

	encryptionPassword, err := HashPassword(password)

	if err != nil {

		test.Fatalf("Failed  case 1 :%v", err)
	}

	if encryptionPassword == password {
		test.Fatalf("Failed case 2:%v", err)
	}

	matchedPassword := CheckPasswordHash(password, encryptionPassword)

	if !matchedPassword {
		test.Fatalf("Failed case 3:%v", err)
	}

	fmt.Println("Success")

}

func TestPasswordHashSet(test *testing.T) { 
	testCases := []struct {
		caseName        string
		correctPassword string
		input           string
		expected        bool
	}{
		{"Correct", "password", "password", true},
		{"Wrong", "password", "wrongPassword", false},
		{"Empty", "password", "", false},
	}

	for _, row := range testCases { 

		test.Run(row.caseName, func(t *testing.T) {

			encryptionPassword, err := HashPassword(row.correctPassword)
			if err != nil {
				t.Fatalf("Setup Failed: %s, err: %v", row.caseName, err)
			}

			matchedPassword := CheckPasswordHash(row.input, encryptionPassword)

			if matchedPassword != row.expected {
				t.Errorf("%s Failed: input '%s', expected %v but got %v",
					row.caseName, row.input, row.expected, matchedPassword)
			}
			t.Logf("Success : %v", row.caseName)
		})

	}
}

func TestRegisterNewUser(test *testing.T) {

	if db == nil {
		db = DatabaseConnection()
	}

	testCases := []struct {
		caseName        string
		shopID          int
		username        string
		correctPassword string
		expected        bool
	}{
		{"Correct_Case", 1, "not_saharat3", "password123", true},
		{"Empty_Username", 1, "", "password123", false},
	}
	for _, row := range testCases {
		test.Run(row.caseName, func(t *testing.T) {

			isSuccess := RegisterNewUser(row.shopID, row.username, row.correctPassword)
			if isSuccess != row.expected {
				t.Errorf("%s Failed: input '%s', expected %v but got %v",
					row.caseName, row.username, row.expected, isSuccess)
			}
			t.Logf("Success : %v", row.caseName)
		})
	}
}
