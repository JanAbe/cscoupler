package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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
