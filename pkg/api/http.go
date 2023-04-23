package api

import (
	"fmt"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

type HTTPServer struct {
	mux *http.ServeMux
}

func NewHTTPServer(workflowManager genconnect.ManagerServiceHandler, projectService genconnect.ProjectServiceHandler) *HTTPServer {
	mux := http.NewServeMux()
	// The generated constructors return a path and a plain net/http
	// handler.
	route, handler := genconnect.NewManagerServiceHandler(workflowManager)
	mux.Handle(route, handler)

	projectRoutes, projectHandlers := genconnect.NewProjectServiceHandler(projectService)

	mux.Handle(projectRoutes, projectHandlers)

	reflector := grpcreflect.NewStaticReflector(
		"workflow.ManagerService",
		"project.ProjectService",
		// protoc-gen-connect-go generates package-level constants
		// for these fully-qualified protobuf service names, so you'd more likely
		// reference userv1.UserServiceName and groupv1.GroupServiceName.
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	// Many tools still expect the older version of the server reflection API, so
	// most servers should mount both handlers.
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	return &HTTPServer{
		mux: mux,
	}
}

func (h *HTTPServer) Serve(port int) error {
	return http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		h2c.NewHandler(h.mux, &http2.Server{}),
	)
}
