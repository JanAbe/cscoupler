package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/janabe/cscoupler/domain"

	h "github.com/janabe/cscoupler/handlers"
)

func main() {
	fmt.Println("Running server on port :3000...")
	registerHandlers()
	user, _ := domain.NewUser("bob@email.com", "password")
	skills := []string{"C#", "linux"}
	exp := []string{"Worked at X", "Interned at Y"}
	student, _ := domain.NewStudent("bob", "fisher", "university of fish", time.Now(), skills, exp, *user, domain.Available)
	fmt.Println(student)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// Function that contains all api paths and registers them
func registerHandlers() {
	http.Handle("/", http.FileServer(http.Dir("./views")))
	http.Handle("/signin", h.LoggingHandler(os.Stdout, h.SigninHandler))
	http.Handle("/signup", h.LoggingHandler(os.Stdout, h.SignupHandler))
	http.Handle("/useronly", h.LoggingHandler(os.Stdout, h.ValidateHandler(h.Useronly)))
}
