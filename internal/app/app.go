package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/caviarman/garm/internal/server"
)

func Run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		dir, file := path.Split(r.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			http.ServeFile(w, r, "./web/dist/web/index.html")
		} else {
			http.ServeFile(w, r, "./web/dist/web/"+path.Join(dir, file))
		}

	})

	httpServer := server.New(r, server.Port("8080"))

	waitSignal(httpServer)

	return nil
}

func waitSignal(httpServer *server.Server) {
	fmt.Println("App started!")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Println("shutdown signal: " + s.String())
	case err := <-httpServer.Notify():
		fmt.Println(err, "waitSignal - httpServer.Notify")
	}

	err := httpServer.Shutdown()
	if err != nil {
		fmt.Println(err, "waitSignal - httpServer.Shutdown")
	}
}
