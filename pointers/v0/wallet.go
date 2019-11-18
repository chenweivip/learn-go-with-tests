package main

import "fmt"

// Wallet stores the number of Bitcoin someone owns
type Wallet struct {
	balance int
}

// Deposit will add some Bitcoin to a wallet
func (w Wallet) Deposit(amount int) {
	fmt.Printf("address of balance in Deposit is %v \n", &w.balance)
	w.balance += amount
}

// Balance returns the number of Bitcoin a wallet has
func (w Wallet) Balance() int {
	return w.balance
}
