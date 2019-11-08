package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/janabe/cscoupler/domain"
	"github.com/janabe/cscoupler/services"
)

// ProjectHandler struct containing all
// project related handler funcs
type ProjectHandler struct {
	ProjectService services.ProjectService
	AuthHandler    AuthHandler
	Path           string
}

// ProjectData is a struct that corresponds to incoming project data
type ProjectData struct {
	ID              string   `json:"id"`
	Description     string   `json:"description"`
	Compensation    string   `json:"compensation"`
	Duration        string   `json:"duration"`
	Recommendations []string `json:"recommendations"`
	CompanyID       string   `json:"companyID"`
}

// FetchAllProjects fetches all projects
func (p ProjectHandler) FetchAllProjects() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "OPTIONS" {
			return
		}

		projects, err := p.ProjectService.FetchAll()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		var projectsData []ProjectData
		for _, p := range projects {
			projectsData = append(projectsData, ToProjectData(p))
		}

		json.NewEncoder(w).Encode(projectsData)
	})
}

// DeleteProject deletes a project
func (p ProjectHandler) DeleteProject() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			return
		}

		// todo: validate that the representative removing the project
		// is an employee of the company that listed the project.
		// Otherwise you can delete projects from other companies.

		projectID := strings.TrimPrefix(r.URL.Path, p.Path+"delete/")
		_, err := p.ProjectService.FindByID(projectID)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = p.ProjectService.Delete(projectID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// Register registers all project related handlers
func (p ProjectHandler) Register() {
	http.Handle(p.Path, LoggingHandler(os.Stdout, p.AuthHandler.Validate("", p.FetchAllProjects())))
	http.Handle(p.Path+"delete/", LoggingHandler(os.Stdout, p.AuthHandler.Validate(domain.RepresentativeRole, p.DeleteProject())))
}
