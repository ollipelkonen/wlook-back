package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ollipelkonen/wlook_back/config"
	"github.com/ollipelkonen/wlook_back/repositories"
	"github.com/ollipelkonen/wlook_back/services"
	ping "github.com/ollipelkonen/wlook_back/target_test"
)

func authMiddleware(token string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.Header["Authorization"]) > 0 {
				if r.Header["Authorization"][0] == ("Bearer " + token) {
					next.ServeHTTP(w, r)
				}
			}
		})
	}
}

func main() {
	config := config.LoadConfig("settings.json")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database)

	todoRepository := repositories.CreateRepository(connectionString)
	todoService := services.CreateService(todoRepository)

	fmt.Printf("____ listening to port %s\n", config.Port)

	r := mux.NewRouter()

	r.Handle("/todo", todoService.CreateGetAllEndpoint()).Methods("GET")
	r.Handle("/todo/{id}", todoService.CreateGetByIdEndpoint()).Methods("GET")
	r.Handle("/todo", todoService.CreateInsertEndpoint()).Methods("POST")
	r.Handle("/todo/{id}", todoService.CreateUpdateEndpoint()).Methods("PUT")
	r.Handle("/todo/{id}", todoService.CreateDeleteEndpoint()).Methods("DELETE")

	r.Use(authMiddleware(config.Token))
	http.ListenAndServe(":"+config.Port, r)

	ping.Ping.Test()
}
