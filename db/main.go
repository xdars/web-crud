package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"github.com/xdars/web-crud/graph/model"
)

type dbm struct {
	*sql.DB
}


func GetDatabase() (*dbm, error) {
	db, err := sql.Open("sqlite3", "db/db.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &dbm{db}, nil
}

func (h *dbm) CreateUser(id, firstname, lastname string) (bool, error) {
	tx, err := h.Begin()
	if err != nil  {
		log.Println(err)
		return false, err
	}

	stmt, err := tx.Prepare("insert into users(id, firstname, lastname) values(?, ?, ?)")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	stmt.Exec(id, firstname, lastname)
	tx.Commit()
	return true, nil
}

func (h *dbm) GetUsers(m *[]*model.User) error {
	rows, err := h.Query("select * from users")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var fn string
		var ln string
		err = rows.Scan(&id, &fn, &ln)
		if err != nil {
			log.Fatal(err)
		}
		*m = append(*m, &model.User{id, fn, ln})
	}
	return nil
}

func (h *dbm) GetUser(id string, m *model.User) error {
	stmt, err := h.Prepare("select * from users where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&m.ID, &m.FirstName, &m.LastName)
	if err != nil {
		return err
	}
	return nil
}