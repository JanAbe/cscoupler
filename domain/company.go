package domain

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// todo: look into which functions i want the different structs to have
// + look into which fields i want to export.
// do i want these fields to be changed by the user however the want
// or do i only allow this via a method of this struct
// e.g add a AddAddress(addr) method or not

// todo: look into which structs need an ID field

// CompanyRepository interface
type CompanyRepository interface {
	Create(company Company) error
	FindAll() ([]Company, error)
	FindByID(id string) (Company, error)
	FindByName(name string) (Company, error)
}

// Company struct conveying a company
// that is looking for skilled students
type Company struct {
	ID              string
	Name            string
	Description     string
	Locations       []Address
	Representatives []Representative
}

// NewCompany creates a new Company based on the
// provided input if all input is valid, returning
// an error otherwise
func NewCompany(name, desc string) (Company, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return Company{}, errors.New("provided name can't be empty")
	}

	if len(strings.TrimSpace(desc)) == 0 {
		return Company{}, errors.New("provided description can't be empty")
	}

	return Company{
		ID:              uuid.New().String(),
		Name:            strings.ToLower(name),
		Description:     strings.ToLower(desc),
		Locations:       []Address{},
		Representatives: []Representative{},
	}, nil
}

// Address struct conveying the addresses
// a company has branches at
type Address struct {
	Street  string
	Zipcode string
	City    string
	Number  string
}

// NewAddress creates a new Addres based on the
// provided input if all input is valid, returning
// an error otherwise
func NewAddress(street, zipcode, city, number string) (Address, error) {
	if len(strings.TrimSpace(street)) == 0 {
		return Address{}, errors.New("provided street can't be empty")
	}

	r := regexp.MustCompile(`^\d{4}\s[A-Z]{2}$`)
	if !r.MatchString(zipcode) {
		return Address{}, errors.New("provided zipcode is invalid, should be of format 0000 XX, where 0 can be any number and X can be any letter")
	}

	if len(strings.TrimSpace(city)) == 0 {
		return Address{}, errors.New("provided city can't be empty")
	}

	if len(strings.TrimSpace(number)) == 0 {
		return Address{}, errors.New("provided number can't be empty")
	}

	return Address{
		Street:  strings.ToLower(street),
		Zipcode: zipcode,
		City:    strings.ToLower(city),
		Number:  strings.ToLower(number),
	}, nil
}

// RepresentativeRepository interface
type RepresentativeRepository interface {
	Create(representative Representative) error
	FindByID(id string) (Representative, error)
}

// Representative struct conveying a
// representative/employee of a company
// that is looking in name of the company
// for students
type Representative struct {
	ID        string
	Position  string
	User      User
	CompanyID string
}

// NewRepresentative creates a new representative based on the provided input
func NewRepresentative(pos, companyID string, user User) (Representative, error) {
	if len(strings.TrimSpace(pos)) == 0 {
		return Representative{}, errors.New("provided position can't be empty")
	}

	return Representative{
		ID:        uuid.New().String(),
		Position:  strings.ToLower(pos),
		User:      user,
		CompanyID: companyID,
	}, nil
}

// Project struct conveying a
// project for which the company
// is looking for students
/*
compensation, duration and recommendations
are strings and not some other datatype
because I want to give the user freedom
in how to express themselves.
e.g duration= 3-4 months | negotiable
*/
type Project struct {
	ID              string
	Description     string
	Compensation    string
	Duration        string
	Recommendations []string
}

// NewProject creates a new Project based on the
// provided input if all is valid, it returns
// an error otherwise
func NewProject(desc, comp, dur string) (Project, error) {
	if len(strings.TrimSpace(desc)) == 0 {
		return Project{}, errors.New("provided description can't be empty")
	}

	if len(strings.TrimSpace(comp)) == 0 {
		return Project{}, errors.New("provided compensation can't be empty")
	}

	if len(strings.TrimSpace(dur)) == 0 {
		return Project{}, errors.New("provided duration can't be empty")
	}

	return Project{
		ID:              uuid.New().String(),
		Description:     strings.ToLower(desc),
		Compensation:    strings.ToLower(comp),
		Duration:        strings.ToLower(dur),
		Recommendations: []string{},
	}, nil
}
