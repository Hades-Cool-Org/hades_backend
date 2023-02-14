package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	customMiddleware "hades_backend/api/middleware"
	"hades_backend/app/database"
	"net/http"
)

var (
	db = database.DB
)

func Handler(l *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(customMiddleware.Logger(l))
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Creating mysql default handler
	mysqlHandler := NewMySQLHandler(db)

	r.Group(mysqlHandler.Handle)

	//utils.GenerateDocs(r)
	return r
}
