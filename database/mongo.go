package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func ConnectMongo() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI("mongodb://localhost:27017"),
	)
	if err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}

	// Testar a conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Erro ao pingar MongoDB:", err)
	}

	Client = client
	log.Println("Conectado ao MongoDB com sucesso")
}
