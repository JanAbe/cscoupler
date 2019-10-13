package handlers

import (
	d "github.com/janabe/cscoupler/domain"
)

// mappers contains functions that map *Data struct to domain structs
// and vice versa

// ToCompanyData maps a company domain struct to
// a companyData struct
func ToCompanyData(c d.Company) CompanyData {
	companyData := CompanyData{
		Name:        c.Name,
		Information: c.Information,
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
			JobTitle:  r.JobTitle,
			CompanyID: r.CompanyID,
			UserData: UserData{
				Email:     r.User.Email,
				Firstname: r.User.FirstName,
				Lastname:  r.User.LastName,
			},
		}

		companyData.Representatives = append(companyData.Representatives, reprData)
	}

	return companyData
}

// ToStudentData maps a student domain struct to
// a studentData struct
func ToStudentData(s d.Student) StudentData {
	studentData := StudentData{
		University: s.University,
		Skills:     s.Skills,
		Experience: s.Experience,
		UserData: UserData{
			Email:     s.User.Email,
			Firstname: s.User.FirstName,
			Lastname:  s.User.LastName,
		},
	}

	return studentData
}

// ToRepresentativeData maps a representative domain struct
// to a representativeData struct
func ToRepresentativeData(r d.Representative) RepresentativeData {
	representativeData := RepresentativeData{
		JobTitle:  r.JobTitle,
		CompanyID: r.CompanyID,
		UserData: UserData{
			Email:     r.User.Email,
			Firstname: r.User.FirstName,
			Lastname:  r.User.LastName,
		},
	}

	return representativeData
}
