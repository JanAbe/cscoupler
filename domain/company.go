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
	Information     string
	Locations       []Address
	Representatives []Representative
	Projects        []Project
}

// NewCompany creates a new Company based on the
// provided input if all input is valid, returning
// an error otherwise
func NewCompany(name, info string) (Company, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return Company{}, errors.New("provided name can't be empty")
	}

	if len(strings.TrimSpace(info)) == 0 {
		return Company{}, errors.New("provided information can't be empty")
	}

	return Company{
		ID:              uuid.New().String(),
		Name:            strings.ToLower(name),
		Information:     strings.ToLower(info),
		Locations:       []Address{},
		Representatives: []Representative{},
		Projects:        []Project{},
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
	CompanyID       string
}

// NewProject creates a new Project based on the
// provided input if all is valid, it returns
// an error otherwise
func NewProject(desc, comp, dur, companyID string) (Project, error) {
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
		CompanyID:       companyID,
	}, nil
}
