package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// URL is the base URL to the webserver
var URL = "https://localhost:3000"

// GetJWTSecret gets the jwt-secret from the provided file
func GetJWTSecret(filepath string) []byte {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var secret jwtSecret
	err = json.Unmarshal(data, &secret)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return []byte(secret.JWTSecret)
}

// Helper struct to unmarshal json into secret
type jwtSecret struct {
	JWTSecret string `json:"jwtsecret"`
}

// GetDSN gets the data-source-name from the provided file
func GetDSN(filepath string) string {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var dsn dsn
	err = json.Unmarshal(data, &dsn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return dsn.DSN
}

type dsn struct {
	DSN string `json:"dsn"`
}
