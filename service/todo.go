package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt, err := s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var TODO model.TODO
	TODO.ID = id
	err = stmt.QueryRowContext(ctx, id).Scan(&TODO.Subject, &TODO.Description, &TODO.CreatedAt, &TODO.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, err

	}
	return &TODO, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
		readAll    = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC`
	)
	TODOs := []*model.TODO{}

	if prevID == 0 {
		if size == 0 {
			rows, err := s.db.Query(readAll)
			switch {
			case err == sql.ErrNoRows:
				return TODOs, err
			case err != nil:
				return TODOs, err
			}

			for rows.Next() {
				var TODO model.TODO
				err = rows.Scan(&TODO.ID, &TODO.Subject, &TODO.Description, &TODO.CreatedAt, &TODO.UpdatedAt)
				if err != nil {
					break
				}
				TODOs = append(TODOs, &TODO)
			}

			return TODOs, nil

		}
		stmt, err := s.db.PrepareContext(ctx, read)
		if err != nil {
			return TODOs, err
		}
		defer stmt.Close()

		rows, err := stmt.QueryContext(ctx, size)
		switch {
		case err == sql.ErrNoRows:
			return TODOs, err
		case err != nil:
			return TODOs, err
		}

		for rows.Next() {
			var TODO model.TODO
			err = rows.Scan(&TODO.ID, &TODO.Subject, &TODO.Description, &TODO.CreatedAt, &TODO.UpdatedAt)
			if err != nil {
				break
			}
			TODOs = append(TODOs, &TODO)
		}

		return TODOs, nil
	} else {
		stmt, err := s.db.PrepareContext(ctx, readWithID)
		if err != nil {
			return TODOs, err
		}
		defer stmt.Close()

		rows, err := stmt.QueryContext(ctx, prevID, size)
		switch {
		case err == sql.ErrNoRows:
			return TODOs, err
		case err != nil:
			return TODOs, err
		}

		for rows.Next() {
			var TODO model.TODO
			err = rows.Scan(&TODO.ID, &TODO.Subject, &TODO.Description, &TODO.CreatedAt, &TODO.UpdatedAt)
			if err != nil {
				break
			}
			TODOs = append(TODOs, &TODO)
		}

		return TODOs, nil

	}
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, &model.ErrNotFound{}
	}

	stmt, err := s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var TODO model.TODO
	TODO.ID = id
	err = stmt.QueryRowContext(ctx, id).Scan(&TODO.Subject, &TODO.Description, &TODO.CreatedAt, &TODO.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, err

	}
	return &TODO, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf(deleteFmt, strings.Repeat(`,?`, len(ids)-1)))
	if err != nil {
		return err
	}
	defer stmt.Close()

	idsArr := make([]interface{}, len(ids))
	for i, v := range ids {
		idsArr[i] = v
	}

	result, err := stmt.ExecContext(ctx, idsArr...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
