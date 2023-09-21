package database

import (
	"context"
	"errors"
	"fmt"
	"main/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ArrUser struct {
	Phone     []uint64
	Name      []string
	Block     []bool
	CreatedAt []time.Time
	UpdatedAt []time.Time
}

type DB struct {
	db *pgx.Conn
}

func NewDatabase() *DB {
	return &DB{}
}

func (ps *DB) ConnectDB(user, pass, host, port, dbname string) error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbname)
	db, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return err
	}
	ps.db = db
	return nil
}

func (ps *DB) GetUsers(arr *ArrUser) error {
	sql_statement := `SELECT phone, name, updated_at FROM users WHERE 
										updated_at BETWEEN NOW() - interval '10 day' and NOW()
										ORDER BY updated_at DESC`
	rows, err := ps.db.Query(context.Background(), sql_statement)
	if err != nil {
		return err
	}
	var u config.User
	for rows.Next() {
		err := rows.Scan(&u.Phone, &u.Name, &u.UpdatedAt)
		if err != nil {
			return err
		}
		arr.Phone = append(arr.Phone, u.Phone)
		arr.Name = append(arr.Name, u.Name)
		arr.UpdatedAt = append(arr.UpdatedAt, u.UpdatedAt)
	}
	return nil
}

func (q *DB) CreateUser(Phone uint64, Name string) error {
	var duplicateEntryError = &pgconn.PgError{Code: "23505"}
	dsn := `INSERT INTO users (phone, name, block, created_at, updated_at) VALUES ($1, $2, $3, Now(), Now())`
	_, err := q.db.Exec(context.Background(), dsn, Phone, Name, false)
	if err != nil {
		if errors.As(err, &duplicateEntryError) {
			fmt.Println("Duplicate key found.")
			fmt.Println(Name, Phone)
			sql := "UPDATE users SET updated_at = $1, name =$2 WHERE phone = $3;"
			_, err := q.db.Exec(context.Background(), sql, time.Now(), Name, Phone)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	fmt.Println(Name, Phone)
	return err
}

func (ps *DB) DeleteUser(Phone uint64) (int64, error) {
	sql := "DELETE FROM users WHERE phone = $1"

	res, err := ps.db.Exec(context.Background(), sql, Phone)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	num := res.RowsAffected()
	return num, err
}
