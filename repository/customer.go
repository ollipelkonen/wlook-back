package repository

/*
	This package contains code to for Customer MySQL repository. Uses MySQL
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

type CustomerRepository interface {
	connect(string)
	GetAll() ([]Customer, error)
	GetById(id string) (Customer, error)
	DeleteById(id string) (string, error)
	Insert(map[string]interface{}) (string, error)
	Update(id string, data map[string]interface{}) (string, error)
	createDatabaseIfNotExists()
}

type CustomerRepositoryImpl struct {
	db *sqlx.DB
}

// initialize repository and connect to database
func CreateCustomerRepository(connectString string) CustomerRepositoryImpl {
	rep := CustomerRepositoryImpl{}
	rep.connect(connectString)
	rep.createDatabaseIfNotExists()
	return rep
}

// create database connection
func (repo *CustomerRepositoryImpl) connect(connectString string) {
	fmt.Println("connecting " + connectString)
	repo.db = sqlx.MustConnect("mysql", connectString)
	//defer repo.db.Close()
}

func (repo *CustomerRepositoryImpl) createDatabaseIfNotExists() {
	var id int
	err := repo.db.Get(&id, "SELECT count(*) FROM customer")
	if err != nil {
		if _, err := repo.db.Exec(customerSchema); err != nil {
			panic(fmt.Sprintf("Customer not found - unable to create\n %s", err))
		} else {
			fmt.Println("Table created")
		}
	}
}

// helper function: try to populate struct Customer with values in map
func mapToCustomer(dict map[string]interface{}) Customer {
	customer := &Customer{}
	t := reflect.ValueOf(customer).Elem()
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
	return *customer
}

// Interface implementation - Mostly self-documenting database functions ahead

func (repo *CustomerRepositoryImpl) GetAll() ([]Customer, error) {
	result := []Customer{}
	err := repo.db.Select(&result, "SELECT * FROM customer")
	if err != nil {
		panic(fmt.Sprintf("!! error: %+v\n", err))
	}
	return result, err
}

func (repo *CustomerRepositoryImpl) GetById(id string) (Customer, error) {
	customer := Customer{}
	err := repo.db.Get(&customer, "SELECT * FROM customer WHERE id = ?", id)
	return customer, err
}

func (repo *CustomerRepositoryImpl) DeleteById(id string) (string, error) {
	_, err := repo.db.Query("DELETE FROM customer WHERE id = ?", id)
	return "ok", err
}

func (repo *CustomerRepositoryImpl) Insert(data map[string]interface{}) (string, error) {
	keys, x := parseMapForQuery(data)
	query := "INSERT INTO customer (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(x, ",") + ")"
	values := mapToCustomer(data)
	_, err := repo.db.NamedExec(query, values)
	return "ok", err
}

func (repo *CustomerRepositoryImpl) Update(id string, data map[string]interface{}) (string, error) {
	// create list of "name=:name, desc=:desc..." for query
	keys, _ := parseMapForQuery(data)
	vals := []string{}
	for _, v := range keys {
		vals = append(vals, v+"=:"+v)
	}
	query := "UPDATE customer SET " + strings.Join(vals, ",") + " WHERE id=:id"
	values := mapToCustomer(data)
	values.Id, _ = strconv.Atoi(id)
	_, err := repo.db.NamedExec(query, values)
	return "ok", err
}
