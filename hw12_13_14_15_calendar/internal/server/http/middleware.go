package internalhttp

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app" //nolint:depguard
)

func loggingMiddleware(logger app.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrw := &WrapResponseWriter{ResponseWriter: w}
		sb := strings.Builder{}
		sb.WriteString(r.RemoteAddr)
		sb.WriteRune(' ')
		sb.WriteString(r.Method)
		sb.WriteRune(' ')
		sb.WriteString(r.RequestURI)
		sb.WriteRune(' ')
		sb.WriteString(r.Proto)
		sb.WriteRune(' ')

		startTime := time.Now()
		next.ServeHTTP(wrw, r)
		duration := time.Since(startTime).Milliseconds()

		sb.WriteString(strconv.FormatInt(int64(wrw.status), 10))
		sb.WriteRune(' ')
		sb.WriteString(strconv.FormatInt(duration, 10))
		sb.WriteRune(' ')
		sb.WriteString(r.Header.Get("User-Agent"))

		logger.Info(sb.String())
	})
}

type WrapResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *WrapResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}
