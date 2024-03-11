package server

import (
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"html/template"
	"net/http"
)

type Studio struct {
	projectService genconnect.ProjectServiceHandler
	defaultProject *gen.Project
}

func NewStudio(
	projectService genconnect.ProjectServiceHandler,
	defaultProject *gen.Project,
) *Studio {
	return &Studio{
		projectService: projectService,
		defaultProject: defaultProject,
	}
}

func (s *Studio) Router() *chi.Mux {
	r := chi.NewRouter()
	logger := httplog.NewLogger("studio", httplog.Options{
		JSON: true,
	})
	r.Use(httplog.RequestLogger(logger))

	// TODO breadchris for each route that will render a whole page, you must specify the pattern to include and then when executing
	// the template, you must specify "layout"
	tmpl := func(w http.ResponseWriter, name string, data any) error {
		return template.Must(template.ParseGlob("pkg/api/html/*.gohtml")).ExecuteTemplate(w, fmt.Sprintf("%s.gohtml", name), data)
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var project *gen.Project
		p, err := s.projectService.GetProject(r.Context(), connect.NewRequest(&gen.GetProjectRequest{}))
		if err != nil {
			project = s.defaultProject
		} else {
			project = p.Msg.Project
		}
		err = tmpl(w, "layout", project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	r.Post("/node", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl(w, "node", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return r
}
