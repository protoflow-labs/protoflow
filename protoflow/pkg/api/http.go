package api

import (
	"github.com/breadchris/protoflow/gen/workflow"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func NewAPIHandler(workflowManager workflow.Manager) http.Handler {
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

	twirpHandler := workflow.NewManagerServer(workflowManager)
	muxRoot.Handle(twirpHandler.PathPrefix(), twirpHandler)
	return muxRoot
}
