package configs

import (
	"context"
	"fmt"
	"log"
	"project-mngt-golang-graphql/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type DB struct {
	client *mongo.Client
}


//fonction qui occupe la connection entre la base des donnes
func ConnectDB() *DB {
	
	// Définir un contexte avec un timeout de 10 secondes
	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second) 
	defer cancel() // S'assurer que le contexte est annulé après usage
	
	// Connexion à MongoDB avec `mongo.Connect`
	client , err := mongo.Connect(ctx , options.Client().ApplyURI(LoadEnvUrl())) //recupere la fonction helper du loading du .env
	
	if err != nil {
		log.Fatal("Il s'agit l'erreur entre la connexion a la base de donnes peut etre le lien : ",err)
	}


	//ping the database
	err = client.Ping(ctx , nil)
	if err != nil {
		log.Fatal("Erreur du ping de database: " , err)
	}

	fmt.Println("\u2707 Connected to mongodb")

	return &DB{client: client}
}

//fonction pour acceder a notre base des donnes et retourne la reference a cette collection
func colHelper(db *DB , collectionName string) *mongo.Collection{
	return db.client.Database("ProjectManagement").Collection(collectionName)
}

//pour creer des projets
func(db *DB) CreateProject(input *model.NewProject) (*model.Project , error) {
	
	//cette variable collection contient notre collection project
	collection := colHelper(db , "Project")
	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()


	res , err := collection.InsertOne(ctx, input)

	if err != nil {
		return nil , err
	}

	project := &model.Project{
		ID: 			res.InsertedID.(primitive.ObjectID).Hex(),
		Owner: 			input.Owner,
		Name: 			input.Name,
		Description: 	input.Description,
		Status: 		model.StatusNotStarted,	
	}

	return project , err
}