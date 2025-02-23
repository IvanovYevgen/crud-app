package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/crud-app/internal/domain"
)

// / We need this to keep the database connection
type BooksDatabase struct {
	db *sql.DB
}

// / Incapsulation of the database connection
func NewBooksDatabase(db *sql.DB) *BooksDatabase {
	return &BooksDatabase{db}
}

func (b *BooksDatabase) CreateBook(ctx context.Context, book domain.Book) error {
	_, err := b.db.Exec("INSERT INTO books (title, author, publish_date, rating) values ($1, $2, $3, $4)",
		book.Title, book.Author, book.PublishDate, book.Rating)

	return err
}

func (b *BooksDatabase) GetByID(ctx context.Context, id int64) (domain.Book, error) {
	var book domain.Book
	err := b.db.QueryRow("SELECT id, title, author, publish_date, rating FROM books WHERE id=$1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
	if err == sql.ErrNoRows {
		return book, domain.ErrBookNotFound
	}

	return book, err
}

func (b *BooksDatabase) GetAll(ctx context.Context) ([]domain.Book, error) {
	rows, err := b.db.Query("SELECT id, title, author, publish_date, rating FROM books")
	if err != nil {
		return nil, err
	}

	books := make([]domain.Book, 0)
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, rows.Err()
}

func (b *BooksDatabase) Delete(ctx context.Context, id int64) error {
	_, err := b.db.Exec("DELETE FROM books WHERE id=$1", id)

	return err
}

func (b *BooksDatabase) Update(ctx context.Context, id int64, inp domain.UpdateBookInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *inp.Title)
		argId++
	}

	if inp.Author != nil {
		setValues = append(setValues, fmt.Sprintf("author=$%d", argId))
		args = append(args, *inp.Author)
		argId++
	}

	if inp.PublishDate != nil {
		setValues = append(setValues, fmt.Sprintf("publish_date=$%d", argId))
		args = append(args, *inp.PublishDate)
		argId++
	}

	if inp.Rating != nil {
		setValues = append(setValues, fmt.Sprintf("rating=$%d", argId))
		args = append(args, *inp.Rating)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE books SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, id)

	_, err := b.db.Exec(query, args...)
	return err
}

func (b *BooksDatabase) CreateTable(ctx context.Context,) error {
	// Create table query
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		publish_date DATE,
		rating INT
	);
	`

	// Execute query
	_, err := b.db.Exec(query)
	return err
}
