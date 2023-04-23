package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/protoflow-labs/protoflow/gen/genconnect"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPServer struct {
	mux *chi.Mux
}

func NewHTTPServer(workflowManager genconnect.ManagerHandler, projectService genconnect.ProjectServiceHandler) *HTTPServer {
	muxRoot := chi.NewRouter()

	muxRoot.Use(middleware.RequestID)
	muxRoot.Use(middleware.RealIP)
	muxRoot.Use(middleware.Logger)
	//muxRoot.Use(session.Sessioner(session.Options{
	//	Provider:           "file",
	//	CookieName:         "session",
	//	FlashEncryptionKey: "SomethingSuperSecretThatShouldChange",
	//}))

	//muxRoot.Use(middleware.Recoverer)
	muxRoot.Use(middleware.Timeout(time.Second * 5))

	route, handler := genconnect.NewManagerHandler(workflowManager)
	muxRoot.Handle(route, handler)

	projectRoutes, projectHandlers := genconnect.NewProjectServiceHandler(projectService)
	muxRoot.Handle(projectRoutes, projectHandlers)

	return &HTTPServer{
		mux: muxRoot,
	}
}

func (h *HTTPServer) Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), h.mux)
}
