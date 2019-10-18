package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
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

// HasCorrectContentType checks if the file's
// content type matches the wanted/expected content type
func HasCorrectContentType(file multipart.File, ct string) bool {
	buffer := make([]byte, 521)
	_, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return false
	}

	contentType := http.DetectContentType(buffer)
	return ct == contentType
}

// Capitalize returns a capitalized copy of the word
func Capitalize(word string) string {
	w := strings.ToUpper(string(word[0]))
	return w + word[1:]
}

// CapitalizeLastWord capitalizes the last word
// of the provided string.
// e.g word="van der bilt"
// result="van der Bilt"
func CapitalizeLastWord(word string) string {
	words := strings.Split(word, " ")
	last := words[len(words)-1]
	return strings.Replace(word, last, Capitalize(last), 1)
}
