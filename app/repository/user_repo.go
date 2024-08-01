package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"todo-echo/app/domain/dao"
)

type UserRepository interface {
	Authenticate(ctx context.Context, user dao.User) (int64, error)
	Insert(ctx context.Context, user dao.User) error
}

type UserRepositoryImpl struct {
	Db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{
		Db: db,
	}
}

func (u UserRepositoryImpl) Authenticate(ctx context.Context, user dao.User) (int64, error) {
	plainPassword := user.Password

	script := "SELECT id, email, password, name, created, active FROM user WHERE email=? AND active = TRUE"
	row := u.Db.QueryRowContext(ctx, script, user.Email)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Name, &user.Created, &user.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, dao.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(user.Password, plainPassword)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, dao.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return user.Id, nil
}

func (u UserRepositoryImpl) Insert(ctx context.Context, user dao.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, 12)

	script := "INSERT INTO user(email, password, name, created) VALUES (?, ?, ?, UTC_TIMESTAMP)"
	_, err = u.Db.ExecContext(ctx, script, user.Email, hashedPassword, user.Name)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "email") {
				return dao.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}
