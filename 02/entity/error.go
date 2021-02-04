package entity

import "errors"

// ErrNotFound not found
var ErrNotFound = errors.New("not found")

// ErrIDFromNotFound idFrom not found
var ErrIDFromNotFound = errors.New("idFrom not found")

// ErrIDToNotFound idFrom not found
var ErrIDToNotFound = errors.New("idTo not found")

// ErrNotEnoughbalances not enough balances
var ErrNotEnoughbalances = errors.New("not enough balances")
