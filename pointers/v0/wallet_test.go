package main

import (
	"fmt"
	"testing"
)

func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()

	fmt.Printf("address of balance in test is %v \n", &wallet.balance)

	expected := 10

	if got != expected {
		t.Errorf("got %d expected %d", got, expected)
	}
}
