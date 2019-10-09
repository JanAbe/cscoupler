package services

import (
	"fmt"

	"github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
)

// CompanyService struct, containing all features
// the app supports regarding companies
type CompanyService struct {
	CompanyRepo domain.CompanyRepository
	UserService UserService
}

// Register registers a new company
func (c CompanyService) Register(company domain.Company) error {

	// !!!!!
	// todo: bekijk waar ik wil controleren of een bedrijf/student
	// zijn naam/email al geascocieerd is met een account
	// Momenteel doe ik dit in de handlers, maar is het niet logischer
	// en beter om dat in de services te doen?
	// moet ook kijken naar studentHandler, userHandler + de services
	// het probleem dat ik alleen heb, is dat als ik de checks hier wil uitvoeren
	// (wat beter is vgm), ik niet weet hoe ik specifieke http.Statussen kan
	// returnen in de handler.
	// nvm found an answer, create custom error values.
	// compare returned errors with these erros to determine what went wrong
	if c.NameAlreadyUsed(company.Name) {
		return e.ErrorCompanyNameAlreadyUsed
	}

	// Only one representative, the main account, is present when creating
	// a company
	if c.UserService.EmailAlreadyUsed(company.Representatives[0].User.Email) {
		return e.ErrorEmailAlreadyUsed
	}

	// ok,but what happens if everything goeds right for the
	// representative, and this gets added.
	// But something goes wrong for the company
	// Then there's a representative in the db, but no company...
	err := c.UserService.Register(company.Representatives[0].User)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = c.CompanyRepo.Create(company)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// NameAlreadyUsed checks if a company name already exists or not
func (c CompanyService) NameAlreadyUsed(name string) bool {
	_, err := c.CompanyRepo.FindByName(name)
	if err != nil {
		return false
	}

	return true
}
