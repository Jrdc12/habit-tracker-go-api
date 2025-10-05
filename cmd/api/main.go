package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jrdc12/habit-tracker-go-api/internal/db"
	approuter "github.com/Jrdc12/habit-tracker-go-api/internal/http"
	"github.com/Jrdc12/habit-tracker-go-api/internal/user"
)

func main() {
	// DB setup
	dsn := "file:habiTrack.db?_busy_timeout=5000&_pragma=journal_mode(WAL)"
	sqldb, closeFn, err := db.OpenSQLite(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer closeFn()

	if err := db.InitSchema(sqldb); err != nil {
		log.Fatalf("init schema: %v", err)
	}

	// Repository + Service
	repo := user.NewSQLiteRepository(sqldb)
	service := user.NewService(repo)

	// HTTP router
	mux := approuter.NewRouter(service)

	srv := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		fmt.Println("Listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
