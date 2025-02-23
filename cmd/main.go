package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GOLANG-NINJA/crud-app/internal/repository/psql"
	"github.com/GOLANG-NINJA/crud-app/internal/service"
	"github.com/GOLANG-NINJA/crud-app/internal/transport/rest"
	"github.com/GOLANG-NINJA/crud-app/pkg/database"
	_ "github.com/lib/pq"
)

func main() {
	// init db
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "12345",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// init repo, service, handler
	booksRepo := psql.NewBooksDatabase(db)
	booksService := service.NewBookManager(booksRepo)
	handler := rest.NewHandler(booksService)

	// init & run server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler.InitRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
