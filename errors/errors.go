package errors

import "errors"

// package containing custom errors, used throughout the application

// ErrorEmailAlreadyUsed ...
var ErrorEmailAlreadyUsed = errors.New("email is already in use and bound to an account")

// ErrorCompanyNameAlreadyUsed ...
var ErrorCompanyNameAlreadyUsed = errors.New("company name is already in use")

// ErrorEntityNotFound ...
var ErrorEntityNotFound = errors.New("entity does not exist with the provided id")
