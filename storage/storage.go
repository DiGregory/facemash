package storage

import (
	"database/sql"
)

type Storage struct {
	DB *sql.DB
}

func Connect(Driver, Source string) (*Storage, error) {
	db, err := sql.Open(Driver, Source)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db,
	}, nil

}
