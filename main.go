package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"desafio-todolist-api/database"
	"desafio-todolist-api/repository"
)

func main() {

	// Conexão com o MongoDB
	database.ConnectMongo()

	// Selecionar database
	db := database.Client.Database("taskdb")

	// Criar repository
	taskRepository := repository.NewTaskRepository(db)

	// Criar router
	r := mux.NewRouter()

	// Rota de teste
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API To Do List rodando"))
	}).Methods("GET")

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	_ = taskRepository // temporário
}
