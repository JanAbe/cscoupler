package handlers

import (
	"strconv"
	"strings"

	"github.com/janabe/cscoupler/util"

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

	for _, p := range c.Projects {
		pData := ProjectData{
			Description:     p.Description,
			Compensation:    p.Compensation,
			Duration:        p.Duration,
			Recommendations: p.Recommendations,
			CompanyID:       p.CompanyID,
		}

		companyData.Projects = append(companyData.Projects, pData)
	}

	return companyData
}

// ToStudentData maps a student domain struct to
// a studentData struct
func ToStudentData(s d.Student) StudentData {
	studentData := StudentData{
		ID:               s.ID,
		University:       s.University,
		Skills:           s.Skills,
		Experiences:      s.Experiences,
		ShortExperiences: s.ShortExperiences,
		Wishes:           s.Wishes,
		Status:           ToStatus(strconv.Itoa(int(uint8(s.Status)))),
		Resume:           s.Resume,
		UserData: UserData{
			Email:     s.User.Email,
			Firstname: strings.Title(s.User.FirstName),
			Lastname:  util.CapitalizeLastWord(s.User.LastName),
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

// ToProjectData maps a project domain struct
// to a projectData struct
func ToProjectData(p d.Project) ProjectData {
	projectData := ProjectData{
		ID:              p.ID,
		Description:     p.Description,
		Compensation:    p.Compensation,
		Duration:        p.Duration,
		Recommendations: p.Recommendations,
		CompanyID:       p.CompanyID,
	}

	return projectData
}

// ToStatus transforms the status number
// to the corresponding string representation
func ToStatus(num string) string {
	if num == "0" {
		return "Available"
	}

	if num == "1" {
		return "Unavailable"
	}

	return num
}
