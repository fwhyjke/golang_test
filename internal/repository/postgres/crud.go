package postgres

import (
	"context"
	"fmt"

	"github.com/fwhyjke/golang_test/internal/repository"
)

func (c *ConnectionPool) Create(ctx context.Context, dto repository.NoteDTO) (repository.Note, error) {
	query := `
	INSERT INTO tasks (title, description, done)
	VALUES ($1, $2, $3)
	RETURNING id, title, description, done;
	`

	row := c.QueryRow(ctx, query, dto.Title, dto.Description, dto.Done)
	var newNote repository.Note
	err := row.Scan(
		&newNote.ID,
		&newNote.Title,
		&newNote.Description,
		&newNote.Done,
	)
	if err != nil {
		return repository.Note{}, fmt.Errorf("scan error: %w", err)
	}
	return newNote, nil
}

func (c *ConnectionPool) GetByID(ctx context.Context, id uint64) (repository.Note, error) {
	return repository.Note{}, fmt.Errorf("TODO")
}
func (c *ConnectionPool) GetAll(ctx context.Context) ([]repository.Note, error) {
	return []repository.Note{}, fmt.Errorf("TODO")
}
func (c *ConnectionPool) Update(ctx context.Context, id uint64, dto repository.NoteDTO) (repository.Note, error) {
	return repository.Note{}, fmt.Errorf("TODO")
}
func (c *ConnectionPool) Delete(ctx context.Context, id uint64) error {
	return fmt.Errorf("TODO")
}

func (c *ConnectionPool) CloseConn() {
	c.Close()
}
