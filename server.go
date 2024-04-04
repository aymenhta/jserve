package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         app.cfg.port,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.log.Info(fmt.Sprintf("started listening to http://localhost%s/", app.cfg.port))

	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (app *application) routes() http.Handler {
	// create a server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.listDBHandler)
	mux.HandleFunc("GET /{table}", app.showTableHandler)
	mux.HandleFunc("POST /{table}", app.addRecordHandler)
	mux.HandleFunc("PUT /{table}/{id}", app.editRecordHandler)
	mux.HandleFunc("DELETE /{table}/{id}", app.deleteRecordHandler)

	// table lookup is only through id
	mux.HandleFunc("GET /{table}/{id}", app.showRecordHandler)

	return app.enableCors(app.logRequest(mux))
}

func (app *application) enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", app.cfg.client)
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.log.Info(fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI()))
		next.ServeHTTP(w, r)
	})
}
