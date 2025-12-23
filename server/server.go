package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go-connect-todo/gen/todo/v1/todov1connect"
	"go-connect-todo/gen/transaction/v1/transactionv1connect"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/rs/cors"
)

func New(todo todov1connect.TodoServiceHandler, txnHandler transactionv1connect.AccountServiceHandler) *http.Server {
	mux := http.NewServeMux()

	path, handler := todov1connect.NewTodoServiceHandler(todo, connect.WithInterceptors(validate.NewInterceptor()))
	tPath, tHandler := transactionv1connect.NewAccountServiceHandler(txnHandler, connect.WithInterceptors(validate.NewInterceptor()))

	mux.Handle(path, handler)
	mux.Handle(tPath, tHandler)
	loggedMux := loggingMiddleware(mux)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	s := http.Server{
		Addr:      ":8080",
		Handler:   withCORS(loggedMux),
		Protocols: p,
	}

	return &s
}

func withCORS(mux http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // dev only
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // allow all headers
		AllowCredentials: false,
	})

	return c.Handler(mux)
}

// loggingMiddleware logs the details of an HTTP request and the time taken to process it.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Record the start time

		// Log request details before passing to the next handler
		slog.Info(fmt.Sprintf("Started Request: Method=%s, URL=%s", r.Method, r.URL.Path))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log the completion and elapsed time after the handler returns
		slog.Info(fmt.Sprintf("Completed Request: Method=%s, URL=%s, Duration=%v", r.Method, r.URL.Path, time.Since(start)))
	})
}
