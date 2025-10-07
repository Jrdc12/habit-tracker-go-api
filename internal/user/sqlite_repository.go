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

func (r *sqliteRepo) UpdatePartial(ctx context.Context, id int, name, email *string, passwordHash *string) (User, error) {
	sets := make([]string, 0, 3)
	args := make([]interface{}, 0, 4)

	if name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *name)
	}
	if email != nil {
		sets = append(sets, "email = ?")
		args = append(args, *email)
	}
	if passwordHash != nil {
		sets = append(sets, "password_hash = ?")
		args = append(args, *passwordHash)
	}

	if len(sets) == 0 {
		return r.ByID(ctx, id)
	}

	q := "Update users SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			return User{}, ErrEmailExists
		}
		return User{}, err
	}
	a, _ := res.RowsAffected()
	if a == 0 {
		return User{}, ErrNotFound
	}
	return r.ByID(ctx, id)
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
