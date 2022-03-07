package repository

/*
	This package contains code to for Todo MySQL repository. Uses MySQL
*/

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // required driver, do not remove even if compiler nags it unused
	"github.com/jmoiron/sqlx"
)

type TodoRepository interface {
	connect(string)
	GetAll() ([]Todo, error)
	GetById(id string) (Todo, error)
	DeleteById(id string) (string, error)
	Insert(map[string]interface{}) (string, error)
	Update(id string, data map[string]interface{}) (string, error)
	createDatabaseIfNotExists()
}

type TodoRepositoryImpl struct {
	db *sqlx.DB
}

// initialize repository and connect to database
func CreateRepository(connectString string) TodoRepositoryImpl {
	rep := TodoRepositoryImpl{}
	rep.connect(connectString)
	rep.createDatabaseIfNotExists()
	return rep
}

// create database connection
func (repo *TodoRepositoryImpl) connect(connectString string) {
	fmt.Println("connecting " + connectString)
	repo.db = sqlx.MustConnect("mysql", connectString)
	//defer repo.db.Close()
}

func (repo *TodoRepositoryImpl) createDatabaseIfNotExists() {
	var id int
	err := repo.db.Get(&id, "SELECT count(*) FROM todo")
	if err != nil {
		if _, err := repo.db.Exec(todoSchema); err != nil {
			panic(fmt.Sprintf("Todo not found - unable to create\n %s", err))
		} else {
			fmt.Println("Table created")
		}
	}
}

// helper function: try to populate struct Todo with values in map
func mapToTodo(dict map[string]interface{}) Todo {
	todo := &Todo{}
	t := reflect.ValueOf(todo).Elem()
	for k, v := range dict {
		fname := strings.Title(k)
		f := t.FieldByName(fname)
		if f.CanSet() {
			val := reflect.ValueOf(v)
			if f.Type().String() == "time.Time" {
				newTime, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v", v))
				newValue := reflect.ValueOf(newTime)
				f.Set(newValue)
			} else {
				f.Set(val.Convert(f.Type()))
			}
		}
	}
	return *todo
}

// helper function: return keys in two arrays for sql query, one with names and one with :names
func parseMapForQuery(data map[string]interface{}) ([]string, []string) {
	keys := []string{}
	x := []string{}
	for k, _ := range data {
		keys = append(keys, k)
		x = append(x, ":"+k)
	}
	return keys, x
}

// Interface implementation - Mostly self-documenting database functions ahead

func (repo *TodoRepositoryImpl) GetAll() ([]Todo, error) {
	result := []Todo{}
	err := repo.db.Select(&result, "SELECT * FROM todo")
	if err != nil {
		panic(fmt.Sprintf("!! error: %+v\n", err))
	}
	return result, err
}

func (repo *TodoRepositoryImpl) GetById(id string) (Todo, error) {
	todo := Todo{}
	err := repo.db.Get(&todo, "SELECT * FROM todo WHERE id = ?", id)
	return todo, err
}

func (repo *TodoRepositoryImpl) DeleteById(id string) (string, error) {
	_, err := repo.db.Query("DELETE FROM todo WHERE id = ?", id)
	return "ok", err
}

func (repo *TodoRepositoryImpl) Insert(data map[string]interface{}) (string, error) {
	keys, x := parseMapForQuery(data)
	query := "INSERT INTO todo (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(x, ",") + ")"
	values := mapToTodo(data)
	_, err := repo.db.NamedExec(query, values)
	return "ok", err
}

func (repo *TodoRepositoryImpl) Update(id string, data map[string]interface{}) (string, error) {
	// create list of "name=:name, desc=:desc..." for query
	keys, _ := parseMapForQuery(data)
	vals := []string{}
	for _, v := range keys {
		vals = append(vals, v+"=:"+v)
	}
	query := "UPDATE todo SET " + strings.Join(vals, ",") + " WHERE id=:id"
	values := mapToTodo(data)
	values.Id, _ = strconv.Atoi(id)
	_, err := repo.db.NamedExec(query, values)
	return "ok", err
}
