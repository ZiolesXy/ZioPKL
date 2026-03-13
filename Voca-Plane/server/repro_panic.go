package main

import (
	"fmt"
)

type MyError struct {
	Msg string
}

func (e *MyError) Error() string {
	// Note: midtrans-go's Error() DOES NOT check for nil receiver on the Error struct itself,
	// but it checks for e.RawError != nil. If e is nil, e.RawError panics.
	return "simulated panic"
}

func GetErrorTypedNil() *MyError {
	return nil
}

func GetErrorInterface() error {
	err := GetErrorTypedNil()
	if err != nil {
		return err
	}
	return nil // Explicitly return literal nil
}

func main() {
	err := GetErrorInterface()
	fmt.Printf("Error interface is nil: %v\n", err == nil)
	
	if err == nil {
		fmt.Println("Success: Error is nil, no panic expected.")
		return
	}
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panicked: %v\n", r)
		}
	}()
	
	fmt.Println("Attempting to call err.Error()...")
	_ = err.Error()
}
