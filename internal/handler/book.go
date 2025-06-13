package handler

import (
	"book_rest_api/internal/models"
	"book_rest_api/internal/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BookHandler struct {
	Repo repository.BookRepository
}

func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{Repo: repo}
}

func (bh *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {

	books, err := bh.Repo.GetAllBooks()
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"Operation": "test",
		"What?":     "Getting all books",
		"Books":     books,
	})

}

func (bh *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error: unable to parse id from str to int", http.StatusBadRequest)
		return
	}

	book, err := bh.Repo.GetBookByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"Operation": "test",
		"What?":     "Get a book by id",
		"Book":      book,
	})

}

func (bh *BookHandler) AddBook(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error: failed to read request body", http.StatusInternalServerError)
		return
	}

	var newBook models.Book
	err = json.Unmarshal(bodyBytes, &newBook)
	if err != nil {
		http.Error(w, "error: unable to unmarshal bodyBytes to Book type", http.StatusInternalServerError)
		return
	}

	err = validateBook(newBook)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}

	err = bh.Repo.Create(&newBook)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"Operation": "test",
		"What?":     "Post Request to add a Book",
		"Status":    "Done",
	})

}

func (bh *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error: unable to parse id from str to int", http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error: unable to read bytes from request body", http.StatusInternalServerError)
		return
	}

	var updatedBook models.Book
	err = json.Unmarshal(bodyBytes, &updatedBook)
	if err != nil {
		http.Error(w, "error: unable to unmarshal bytes from body to Book", http.StatusInternalServerError)
		return
	}
	err = validateBook(updatedBook)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}
	updatedBook.ID = id

	err = bh.Repo.Update(&updatedBook)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"Operation": "test",
		"What?":     "Put Request to update a book using id",
		"Status":    "Done",
	})

}

func (bh *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error: unable to parse id from str to int", http.StatusBadRequest)
		return
	}

	err = bh.Repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"Operation": "test",
		"What?":     "Delete request to delete a book by id",
		"Status":    "Done",
	})

}

func validateBook(book models.Book) error {

	if book.Author == "" {
		return fmt.Errorf("book's author can't be empty")
	}

	if book.Title == "" {
		return fmt.Errorf("book's title can't be empty")
	}

	return nil
}
