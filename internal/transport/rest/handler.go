package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	_ "github.com/crud-app/docs"
	"github.com/crud-app/internal/domain"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// BooksManager defines the interface for book management operations.
type BooksManager interface {
	CreateBook(ctx context.Context, book domain.Book) error
	GetByID(ctx context.Context, id int64) (domain.Book, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateBookInput) error
}

type Handler struct {
	booksManager BooksManager
}

func NewHandler(booksManager BooksManager) *Handler {
	return &Handler{
		booksManager: booksManager,
	}
}

func (h *Handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookByID",
			"problem": "invalid ID",
		}).Error(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	book, err := h.booksManager.GetByID(context.TODO(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookByID",
			"problem": "book not found",
		}).Error(err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "unmarshaling request",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.booksManager.CreateBook(context.TODO(), book); err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "reading request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"problem": "invalid ID",
		}).Error(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.booksManager.Delete(context.TODO(), id); err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"problem": "failed to delete book",
		}).Error(err)
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.booksManager.GetAll(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getAllBooks",
			"problem": "failed to fetch books",
		}).Error(err)
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "invalid ID",
		}).Error(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var inp domain.UpdateBookInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "invalid input",
		}).Error(err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.booksManager.Update(context.TODO(), id, inp); err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "failed to update book",
		}).Error(err)
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid ID")
	}
	return id, nil
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	books := r.PathPrefix("/books").Subrouter()
	{
		books.HandleFunc("", h.createBook).Methods(http.MethodPost)
		books.HandleFunc("", h.getAllBooks).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.getBookByID).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.deleteBook).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", h.updateBook).Methods(http.MethodPut)
	}

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}