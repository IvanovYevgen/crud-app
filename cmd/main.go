package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/crud-app/internal/config"
	"github.com/crud-app/internal/repository/psql"
	"github.com/crud-app/internal/service"
	"github.com/crud-app/internal/transport/rest"
	"github.com/crud-app/pkg/database"
	"github.com/crud-app/pkg/hash"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	// init db
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// init deps
	hasher := hash.NewSHA1Hasher("salt")

	booksRepo := psql.NewBooks(db)
	booksService := service.NewBookManager(booksRepo)

	usersRepo := psql.NewUsers(db)
	usersService := service.NewUsers(usersRepo, hasher, []byte("sample secret"), cfg.Auth.TokenTTL)

	handler := rest.NewHandler(booksService, usersService)

	// init & run server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Info("SERVER STARTED")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
