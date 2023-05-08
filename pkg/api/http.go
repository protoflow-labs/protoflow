package api

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/studio/public"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/protoflow-labs/protoflow/gen/genconnect"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type HTTPServer struct {
	config Config
	mux    *http.ServeMux
}

var ProviderSet = wire.NewSet(
	NewConfig,
	NewHTTPServer,
)

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

func NewHTTPServer(
	config Config,
	projectService genconnect.ProjectServiceHandler,
	generateService genconnect.GenerateServiceHandler,
) (*HTTPServer, error) {
	apiMux := http.NewServeMux()

	interceptors := connect.WithInterceptors(NewLogInterceptor())

	// The generated constructors return a path and a plain net/http
	// handler.
	projectRoutes, projectHandlers := genconnect.NewProjectServiceHandler(projectService, interceptors)
	apiMux.Handle(projectRoutes, projectHandlers)

	generateRoutes, generateHandlers := genconnect.NewGenerateServiceHandler(generateService, interceptors)
	apiMux.Handle(generateRoutes, generateHandlers)

	recoverCall := func(_ context.Context, spec connect.Spec, _ http.Header, p any) error {
		log.Error().Msgf("%+v\n", p)
		if err, ok := p.(error); ok {
			return err
		}
		return fmt.Errorf("panic: %v", p)
	}

	reflector := grpcreflect.NewStaticReflector(
		"project.ProjectService",
		"generate.GenerateService",
		// protoc-gen-connect-go generates package-level constants
		// for these fully-qualified protobuf service names, so you'd more likely
		// reference userv1.UserServiceName and groupv1.GroupServiceName.
	)
	apiMux.Handle(grpcreflect.NewHandlerV1(reflector, connect.WithRecover(recoverCall)))
	// Many tools still expect the older version of the server reflection API, so
	// most servers should mount both handlers.
	apiMux.Handle(grpcreflect.NewHandlerV1Alpha(reflector, connect.WithRecover(recoverCall)))

	assets := public.Assets
	fs := http.FS(public.Assets)
	httpFileServer := http.FileServer(fs)

	// TODO breadchris break this up into a separate function
	u, err := url.Parse(config.StudioProxy)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msgf("request: %s", r.URL.Path)

		// If the path is '/studio', forward the request to the other mux handlers
		if r.URL.Path == "/studio" || strings.HasPrefix(r.URL.Path, "/studio/") || r.URL.Path == "/esbuild" {
			r.URL.Path = strings.Replace(r.URL.Path, "/studio", "", 1)

			if config.StudioProxy != "" {
				proxy.ServeHTTP(w, r)
			} else {
				filePath := r.URL.Path
				if strings.Index(r.URL.Path, "/") == 0 {
					filePath = r.URL.Path[1:]
				}

				f, err := assets.Open(filePath)
				if os.IsNotExist(err) {
					r.URL.Path = "/"
				}
				if err == nil {
					f.Close()
				}
				log.Debug().Msgf("serving file: %s", filePath)
				httpFileServer.ServeHTTP(w, r)
			}
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			r.URL.Path = strings.Replace(r.URL.Path, "/api", "", 1)
		}
		apiMux.ServeHTTP(w, r)
		return
	})

	return &HTTPServer{
		config: config,
		mux:    mux,
	}, nil
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
