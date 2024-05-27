package customErrors

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance error")
	ErrAccountNotFound = errors.New("account not found error")
	ErrWrongPassword = errors.New("wrong pin error")
)