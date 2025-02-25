package service

import (
	"context"
	"github.com/crud-app/internal/domain"
	"time"
)

type BooksRepository interface {
	CreateBook(ctx context.Context, book domain.Book) error
	GetByID(ctx context.Context, id int64) (domain.Book, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateBookInput) error
}

type BooksService struct {
	repo BooksRepository
}

func NewBookManager(repo BooksRepository) *BooksService {
	return &BooksService{
		repo: repo,
	}
}

func (b *BooksService) Create(ctx context.Context, book domain.Book) error {
	if book.PublishDate.IsZero() {
		book.PublishDate = time.Now()
	}

	return b.repo.CreateBook(ctx, book)
}

func (b *BooksService) GetByID(ctx context.Context, id int64) (domain.Book, error) {
	return b.repo.GetByID(ctx, id)
}

func (b *BooksService) GetAll(ctx context.Context) ([]domain.Book, error) {
	return b.repo.GetAll(ctx)
}

func (b *BooksService) Delete(ctx context.Context, id int64) error {
	return b.repo.Delete(ctx, id)
}

func (b *BooksService) Update(ctx context.Context, id int64, inp domain.UpdateBookInput) error {
	return b.repo.Update(ctx, id, inp)
}
