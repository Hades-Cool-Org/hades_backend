package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"hades_backend/app/logging"
	"net/http"
	"time"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ww.Header().Set("X-Request-Id", middleware.GetReqID(r.Context()))

			t1 := time.Now().Local()
			defer func() {
				l.Info("Served",
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("lat", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("reqId", middleware.GetReqID(r.Context())))
			}()

			ll := l.With(
				zap.String("reqId", middleware.GetReqID(r.Context())),
				zap.String("proto", r.Proto),
				zap.String("path", r.URL.Path),
			)

			r = r.WithContext(logging.WithLogger(r.Context(), ll))

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
