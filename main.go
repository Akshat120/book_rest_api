package main

import (
	"book_rest_api/internal/config"
	"book_rest_api/internal/handler"
	"book_rest_api/internal/middleware"
	"book_rest_api/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	appConfig, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db, err := sql.Open(appConfig.Database.Type, appConfig.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bookRepo := repository.NewBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)

	r := mux.NewRouter()

	r.Use(middleware.RateLimitMiddleware)

	r.HandleFunc("/health", handler.HealthHandler)

	booksRoute := r.PathPrefix("/books").Subrouter()
	// unprotected routes
	booksRoute.HandleFunc("", bookHandler.GetBooks).Methods("GET")
	booksRoute.HandleFunc("/{id:[0-9]+}", bookHandler.GetBookByID).Methods("GET")

	// protected routes
	protectedRoutes := r.PathPrefix("/books").Subrouter()
	protectedRoutes.Use(middleware.BasicAuthMiddleware)
	protectedRoutes.HandleFunc("", bookHandler.AddBook).Methods("POST")
	protectedRoutes.HandleFunc("/{id:[0-9]+}", bookHandler.UpdateBook).Methods("PUT")
	protectedRoutes.HandleFunc("/{id:[0-9]+}", bookHandler.DeleteBook).Methods("DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         appConfig.Server.Addr,
		WriteTimeout: time.Duration(appConfig.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(appConfig.Server.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(appConfig.Server.IdleTimeout) * time.Second,
	}

	go func() {
		fmt.Printf("App: %s running on port %d\n", appConfig.AppName, appConfig.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServer(): %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

}
