package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)



type DB struct {
	client *mongo.Client
}


//fonction qui occupe la connection entre la base des donnes
func ConnectDB() *DB {
	
	client , err := mongo.NewClient(options.Client().ApplyURI(LoadEnvUrl())) //recupere le fonction helper du loading du .env

	if err != nil {
		log.Fatal("Il s'agit l'erreur entre la connexion a la base de donnes peut etre le lien : ",err)
	}

	ctx , _ := context.WithTimeout(context.Background() , 10*time.Second) //initiate timeout 10seconds
	err = client.Connect(ctx)
	if err != nil { 
		log.Fatal("Il s'agit d'une erreur au temps d'essaie de connexion normalement 10 secodes: " , err)
	}

	//ping the database
	err = client.Ping(ctx , nil)
	if err != nil {
		log.Fatal("Erreur du ping de database: " , err)
	}

	fmt.Println("\u2707 Connected to mongodb")

	return &DB{client: client}
}