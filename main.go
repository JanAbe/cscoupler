package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/janabe/cscoupler/handlers"
)

// HandlerFunc creates a handler from a normal func

func main() {
	fmt.Println("Running server, listening on port 3000...")
	registerHandlers()
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// Function that contains all api paths and registers them
func registerHandlers() {
	http.Handle("/", http.FileServer(http.Dir("./views")))
	http.Handle("/signin", h.LoggingHandler(os.Stdout, h.SigninHandler))
	http.Handle("/signup", h.LoggingHandler(os.Stdout, h.SignupHandler))
	http.Handle("/signup/student", h.LoggingHandler(os.Stdout, h.SignupStudentHandler))
	http.Handle("/students/", h.LoggingHandler(os.Stdout, h.ValidateHandler(h.FetchStudentByID)))
	http.Handle("/useronly", h.LoggingHandler(os.Stdout, h.ValidateHandler(h.Useronly)))
}
