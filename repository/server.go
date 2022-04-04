package repository

/*
	This package contains code to for Server MySQL repository. Uses MySQL
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

type ServerRepository interface {
	connect(string)
	GetAll() ([]Server, error)
	GetById(id string) (Server, error)
	DeleteById(id string) (string, error)
	Insert(map[string]interface{}) (string, error)
	Update(id string, data map[string]interface{}) (string, error)
	createDatabaseIfNotExists()
}

type ServerRepositoryImpl struct {
	db *sqlx.DB
}

// initialize repository and connect to database
func CreateServerRepository(connectString string) ServerRepositoryImpl {
	rep := ServerRepositoryImpl{}
	rep.connect(connectString)
	rep.createDatabaseIfNotExists()
	return rep
}

// create database connection
func (repo *ServerRepositoryImpl) connect(connectString string) {
	fmt.Println("connecting " + connectString)
	repo.db = sqlx.MustConnect("mysql", connectString)
	//defer repo.db.Close()
}

func (repo *ServerRepositoryImpl) createDatabaseIfNotExists() {
	var id int
	err := repo.db.Get(&id, "SELECT count(*) FROM server")
	if err != nil {
		if _, err := repo.db.Exec(serverSchema); err != nil {
			panic(fmt.Sprintf("Server not found - unable to create\n %s", err))
		} else {
			fmt.Println("Table created")
		}
	}
}

// helper function: try to populate struct Server with values in map
func mapToServer(dict map[string]interface{}) Server {
	server := &Server{}
	t := reflect.ValueOf(server).Elem()
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
	return *server
}

// Interface implementation - Mostly self-documenting database functions ahead

func (repo *ServerRepositoryImpl) GetAll() ([]Server, error) {
	result := []Server{}
	err := repo.db.Select(&result, "SELECT * FROM server")
	if err != nil {
		panic(fmt.Sprintf("!! error: %+v\n", err))
	}
	return result, err
}

func (repo *ServerRepositoryImpl) GetById(id string) (Server, error) {
	server := Server{}
	err := repo.db.Get(&server, "SELECT * FROM server WHERE id = ?", id)
	return server, err
}

func (repo *ServerRepositoryImpl) DeleteById(id string) (string, error) {
	_, err := repo.db.Query("DELETE FROM server WHERE id = ?", id)
	return "ok", err
}

func (repo *ServerRepositoryImpl) Insert(data map[string]interface{}) (string, error) {
	keys, x := parseMapForQuery(data)
	query := "INSERT INTO server (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(x, ",") + ")"
	values := mapToServer(data)
	_, err := repo.db.NamedExec(query, values)
	return "ok", err
}

func (repo *ServerRepositoryImpl) Update(id string, data map[string]interface{}) (string, error) {
	// create list of "name=:name, desc=:desc..." for query
	keys, _ := parseMapForQuery(data)
	vals := []string{}
	for _, v := range keys {
		vals = append(vals, v+"=:"+v)
	}
	query := "UPDATE server SET " + strings.Join(vals, ",") + " WHERE id=:id"
	values := mapToServer(data)
	values.Id, _ = strconv.Atoi(id)
	_, err := repo.db.NamedExec(query, values)
	return "ok", err
}
