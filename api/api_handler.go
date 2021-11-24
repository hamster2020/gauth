package api

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/hamster2020/gauth"
)

type apiHandler struct {
	mux   *http.ServeMux
	cfg   gauth.Config
	token gauth.Token
	logic gauth.Logic
}

func NewAPIHandler(cfg gauth.Config, token gauth.Token, logic gauth.Logic) apiHandler {
	h := apiHandler{
		mux:   http.NewServeMux(),
		cfg:   cfg,
		token: token,
		logic: logic,
	}

	h.addRoutes()
	return h
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	crw := &codeResponseWriter{ResponseWriter: w}

	func() {
		defer func() {
			if err := recover(); err != nil {
				errIface, ok := err.(error)
				if !ok {
					errIface = fmt.Errorf("%s", err)
				}

				log.Printf("PANIC: %s\n%s", err, debug.Stack())
				writeJSONError(crw, errIface, http.StatusInternalServerError)
			}
		}()

		if err := loadCaller(h.token, crw, req); err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		h.mux.ServeHTTP(crw, req)
	}()

	duration := time.Since(start)
	caller := getCallerFromRequest(req)
	log.Printf("HTTP: caller=%s method=%s url=%s duration=%s status=%d remote-addr=%s",
		caller.Info(),
		req.Method,
		req.URL.String(),
		duration,
		crw.StatusCode,
		req.RemoteAddr,
	)
}

func (h apiHandler) ListenAndServe() error {
	server := http.Server{
		Addr:         h.cfg.GauthAddress,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      h,
	}

	log.Printf("gauth listening on %s", h.cfg.GauthAddress)
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
