package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	_ "desafio-todolist-api/models"
	"desafio-todolist-api/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// @Summary Criar tarefa
// @Description Cria uma nova tarefa
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body services.CreateTaskInput true "Dados da tarefa"
// @Success 201 {object} models.Task
// @Router /tasks [post]
// POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input services.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	task, err := h.service.Create(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

// @Summary Buscar tarefa por ID
// @Tags tasks
// @Produce json
// @Param id path string true "ID da tarefa"
// @Success 200 {object} models.Task
// @Router /tasks/{id} [get]
// GET /tasks/{id}
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	task, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			writeError(w, http.StatusNotFound, "tarefa não encontrada")
			return
		}
		writeError(w, http.StatusInternalServerError, "erro interno")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// @Summary Listar tarefas
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Router /tasks [get]
// GET (filtro flexível)/tasks?status=pending&priority=high
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	tasks, err := h.service.List(status, priority)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

// @Summary Atualizar tarefa
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID da tarefa"
// @Param task body services.UpdateTaskInput true "Dados da tarefa"
// @Success 200 {object} map[string]string
// @Router /tasks/{id} [put]
// PUT /tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var input services.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			writeError(w, http.StatusNotFound, "tarefa não encontrada")
			return
		}
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "tarefa atualizada com sucesso"})
}

// @Summary Deletar tarefa
// @Tags tasks
// @Param id path string true "ID da tarefa"
// @Success 204
// @Router /tasks/{id} [delete]
// DELETE /tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.service.Delete(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			writeError(w, http.StatusNotFound, "tarefa não encontrada")
			return
		}
		writeError(w, http.StatusInternalServerError, "erro interno")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// tratamento de erros

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidTitle),
		errors.Is(err, services.ErrInvalidStatus),
		errors.Is(err, services.ErrInvalidPriority),
		errors.Is(err, services.ErrDueDateInPast),
		errors.Is(err, services.ErrInvalidDueDate),
		errors.Is(err, services.ErrCompletedCantEdit):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "erro interno")
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
