package main

import (
	"book_rest_api/handler"
	"book_rest_api/repository"
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

	db, err := sql.Open(config.Database.Type, config.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bookRepo := repository.NewBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)

	r := mux.NewRouter()

	r.HandleFunc("/health", handler.HealthHandler)

	booksRoute := r.PathPrefix("/books").Subrouter()
	// unprotected routes
	booksRoute.HandleFunc("", bookHandler.GetBooks).Methods("GET")
	booksRoute.HandleFunc("/{id:[0-9]+}", bookHandler.GetBookByID).Methods("GET")

	// protected routes
	protectedRoutes := r.PathPrefix("/books").Subrouter()
	protectedRoutes.Use(BasicAuthMiddleware)
	protectedRoutes.HandleFunc("", bookHandler.AddBook).Methods("POST")
	protectedRoutes.HandleFunc("/{id:[0-9]+}", bookHandler.UpdateBook).Methods("PUT")
	protectedRoutes.HandleFunc("/{id:[0-9]+}", bookHandler.DeleteBook).Methods("DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         config.Server.Addr,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Server.IdleTimeout) * time.Second,
	}

	go func() {
		fmt.Printf("App: %s running on port %d\n", config.AppName, config.Server.Port)
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
