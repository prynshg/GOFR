package datastore

import (
	"database/sql"
	"fmt"

	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"

	"github.com/prynshg/gooooo/model"

	_ "github.com/go-sql-driver/mysql"
)

type student struct {
	db *sql.DB
}

var (
	mysqlDSN = "root:Pri#@g09$@tcp(localhost:3306)/Students"
	mysqlDB  *sql.DB
)

func New() (*student, error) {
	db, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MySQL!")

	mysqlDB = db
	return &student{db: db}, nil
}

func (s *student) GetByID(ctx *gofr.Context, id string) (*model.Student, error) {
	var resp model.Student

	err := ctx.DB().QueryRowContext(ctx, " SELECT id,name,age,class FROM students where id=$1", id).
		Scan(&resp.ID, &resp.Name, &resp.Age, &resp.Class)
	switch err {
	case sql.ErrNoRows:
		return &model.Student{}, errors.EntityNotFound{Entity: "student", ID: id}
	case nil:
		return &resp, nil
	default:
		return &model.Student{}, err
	}
}

func (s *student) Create(ctx *gofr.Context, student *model.Student) (*model.Student, error) {
	var resp model.Student

	err := ctx.DB().QueryRowContext(ctx, "INSERT INTO students (name, age, class) VALUES($1,$2,$3)"+
		" RETURNING  id,name,age,class", student.Name, student.Age, student.Class).Scan(
		&resp.ID, &resp.Name, &resp.Age, &resp.Class)

	if err != nil {
		return &model.Student{}, errors.DB{Err: err}
	}

	return &resp, nil
}

func (s *student) Update(ctx *gofr.Context, student *model.Student) (*model.Student, error) {
	_, err := ctx.DB().ExecContext(ctx, "UPDATE students SET name=$1,age=$2,class=$3 WHERE id=$4",
		student.Name, student.Age, student.Class, student.ID)
	if err != nil {
		return &model.Student{}, errors.DB{Err: err}
	}

	return student, nil
}

func (s *student) Delete(ctx *gofr.Context, id int) error {
	_, err := ctx.DB().ExecContext(ctx, "DELETE FROM students where id=$1", id)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

func (s *student) Close() {
	s.db.Close()
}
