package handlers

import (
	d "github.com/janabe/cscoupler/domain"
)

// mappers contains functions that map *Data struct to domain structs
// and vice versa

// ToCompanyData maps a company domain struct to
// the companyData struct
func ToCompanyData(c d.Company) CompanyData {
	companyData := CompanyData{
		Name:        c.Name,
		Description: c.Description,
	}

	for _, l := range c.Locations {
		locationData := LocationData{
			Street:  l.Street,
			Zipcode: l.Zipcode,
			City:    l.City,
			Number:  l.Number,
		}

		companyData.Locations = append(companyData.Locations, locationData)
	}

	for _, r := range c.Representatives {
		reprData := RepresentativeData{
			Position: r.Position,
			UserData: UserData{
				Email:     r.User.Email,
				Firstname: r.User.Firstname,
				Lastname:  r.User.Lastname,
			},
		}

		companyData.Representatives = append(companyData.Representatives, reprData)
	}

	return companyData
}
