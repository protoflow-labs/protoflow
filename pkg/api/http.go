package api

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/protoflow-labs/protoflow/gen/genconnect"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type HTTPServer struct {
	mux *http.ServeMux
}

func NewLogInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err != nil {
				log.Error().Msgf("connect error: %+v\n", err)
			}
			return resp, err
		}
	}
	return interceptor
}

func NewHTTPServer(projectService genconnect.ProjectServiceHandler, generateService genconnect.GenerateServiceHandler) *HTTPServer {
	mux := http.NewServeMux()

	interceptors := connect.WithInterceptors(NewLogInterceptor())

	// The generated constructors return a path and a plain net/http
	// handler.
	projectRoutes, projectHandlers := genconnect.NewProjectServiceHandler(projectService, interceptors)
	mux.Handle(projectRoutes, projectHandlers)

	generateRoutes, generateHandlers := genconnect.NewGenerateServiceHandler(generateService, interceptors)
	mux.Handle(generateRoutes, generateHandlers)

	reflector := grpcreflect.NewStaticReflector(
		"project.ProjectService",
		"generate.GenerateService",
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
