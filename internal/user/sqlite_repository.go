package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type sqliteRepo struct{ db *sql.DB }

func NewSQLiteRepository(db *sql.DB) Repository { return &sqliteRepo{db: db} }

func (r *sqliteRepo) Create(ctx context.Context, name, email, passwordHash string) (User, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`, name, email, passwordHash)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			return User{}, ErrEmailExists
		}
		return User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return User{}, err
	}
	return r.ByID(ctx, int(id))
}

func (r *sqliteRepo) ByID(ctx context.Context, id int) (User, error) {
	var u User
	var created string
	row := r.db.QueryRowContext(ctx, `SELECT id, name, email, created_at FROM users WHERE id = ?`, id)
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return u, nil
}

func (r *sqliteRepo) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}
	a, _ := res.RowsAffected()
	if a == 0 {
		return ErrNotFound
	}
	return nil
}
