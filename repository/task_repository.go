package repository

import (
	"context"

	"desafio-todolist-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		collection: db.Collection("tasks"),
	}
}

// Cria task
func (r *TaskRepository) Create(task *models.Task) error {
	_, err := r.collection.InsertOne(context.Background(), task)
	return err
}

// Busca todas as tasks com filtros opcionais
func (r *TaskRepository) FindAll(status, priority string) ([]models.Task, error) {

	filter := bson.M{}

	if status != "" {
		filter["status"] = status
	}

	if priority != "" {
		filter["priority"] = priority
	}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task

	err = cursor.All(context.Background(), &tasks)
	return tasks, err
}

// Busca task por ID
func (r *TaskRepository) FindByID(id string) (*models.Task, error) {

	filter := bson.M{"_id": id}

	var task models.Task

	err := r.collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// Atualiza task
func (r *TaskRepository) Update(id string, update bson.M) error {
	filter := bson.M{"_id": id}

	res, err := r.collection.UpdateOne(
		context.Background(),
		filter,
		bson.M{"$set": update},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Deleta task
func (r *TaskRepository) Delete(id string) error {
	filter := bson.M{"_id": id}

	res, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
