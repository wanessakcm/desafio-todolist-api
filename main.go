package main

import (
	"log"
	"net/http"

	"desafio-todolist-api/database"
	"desafio-todolist-api/handlers"
	"desafio-todolist-api/repository"
	"desafio-todolist-api/services"

	"github.com/gorilla/mux"
)

func main() {
	// Conexão com o MongoDB
	database.ConnectMongo()

	// Seleção do database
	db := database.Client.Database("taskdb")

	// Camadas
	taskRepository := repository.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepository)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Router
	r := mux.NewRouter()

	// Rota de teste
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API To Do List rodando"))
	}).Methods("GET")

	// Rotas do CRUD
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.GetTaskByID).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
