package api

import (
	"fmt"
	"net/http"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type HTTPServer struct {
	mux *http.ServeMux
}

func NewHTTPServer(workflowManager genconnect.ManagerHandler, projectService genconnect.ProjectServiceHandler) *HTTPServer {
	// mux := chi.NewRouter()

	// mux.Use(middleware.RequestID)
	// mux.Use(middleware.RealIP)
	// mux.Use(middleware.Logger)
	// // muxRoot.Use(session.Sessioner(session.Options{
	// // 	Provider:           "file",
	// // 	CookieName:         "session",
	// // 	FlashEncryptionKey: "SomethingSuperSecretThatShouldChange",
	// // }))

	// //muxRoot.Use(middleware.Recoverer)
	// mux.Use(middleware.Timeout(time.Second * 5))

	// route, handler := genconnect.NewManagerHandler(workflowManager)
	// mux.Handle(route, handler)

	// projectRoutes, projectHandlers := genconnect.NewProjectServiceHandler(projectService)

	// mux.Handle(projectRoutes, projectHandlers)

	// reflector := grpcreflect.NewStaticReflector(
	// 	"workflow.Manager",
	// 	"project.ProjectService",
	// 	// protoc-gen-connect-go generates package-level constants
	// 	// for these fully-qualified protobuf service names, so you'd more likely
	// 	// reference userv1.UserServiceName and groupv1.GroupServiceName.
	// )
	// mux.Handle(grpcreflect.NewHandlerV1(reflector))
	// // Many tools still expect the older version of the server reflection API, so
	// // most servers should mount both handlers.
	// mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	mux := http.NewServeMux()

	managerRoute, managerHandlers := genconnect.NewManagerHandler(workflowManager)
	mux.Handle(managerRoute, managerHandlers)

	projectRoutes, projectHandlers := genconnect.NewProjectServiceHandler(projectService)
	mux.Handle(projectRoutes, projectHandlers)

	reflector := grpcreflect.NewStaticReflector(
		"workflow.Manager",
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
		h2c.NewHandler(corsMiddleware(h.mux), &http2.Server{}),
	)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, connect-protocol-version")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
