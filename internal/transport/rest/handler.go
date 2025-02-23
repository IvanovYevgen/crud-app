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

// @Summary Get a book by ID
// @Description Get book details by its ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path int true "Book ID"
// @Success 200 {object} domain.Book
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [get]
func (h *Handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	book, err := h.booksManager.GetByID(context.TODO(), id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// @Summary Create a new book
// @Description Create a new book record
// @Tags books
// @Accept  json
// @Produce  json
// @Param book body domain.Book true "Book object"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Invalid input"
// @Router /books [post]
func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "unmarshaling request",
		}).Error(err)
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

// @Summary Delete a book by ID
// @Description Delete a book record by its ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path int true "Book ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [delete]
func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.booksManager.Delete(context.TODO(), id); err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get all books
// @Description Get a list of all books
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Book
// @Failure 500 {string} string "Internal server error"
// @Router /books [get]
func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.booksManager.GetAll(context.TODO())
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// @Summary Update a book
// @Description Update a book record by its ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path int true "Book ID"
// @Param book body domain.UpdateBookInput true "Updated book object"
// @Success 200 {string} string "Updated"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [put]
func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var inp domain.UpdateBookInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.booksManager.Update(context.TODO(), id, inp); err != nil {
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

// InitRouter initializes the router for mux
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

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
