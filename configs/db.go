package configs

import (
	"context"
	"fmt"
	"log"
	"project-mngt-golang-graphql/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (db *DB) CreateOwner(input *model.NewOwner) (*model.Owner , error){
	collection := colHelper(db , "Owner")

	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()

	res , err := collection.InsertOne(ctx , input)

	if err != nil {
		return nil, err
	}

	owner := &model.Owner{
		ID: 		res.InsertedID.(primitive.ObjectID).Hex() ,
		Name: 		input.Name,
		Email: 		input.Email,
		Phone: 		input.Phone,		
	}

	return owner , err
}


func (db *DB) GetOwners() ([]*model.Owner , error){
	collection := colHelper(db , "Owner")

	ctx,cancel := context.WithTimeout(context.Background() , 10*time.Second)
	var owners []*model.Owner
	defer cancel()

	res , err := collection.Find(ctx , bson.M{})

	if err != nil {
		return nil , err
	}
	defer res.Close(ctx)
	
	for res.Next(ctx) {
		var singleOwner *model.Owner
		if err = res.Decode(&singleOwner); err != nil {
			log.Fatal(err)
		}
		owners = append(owners , singleOwner)
	}

	return owners , err
}


func (db *DB)  GetProject() ([]*model.Project , error){
	collection := colHelper(db , "Project")

	ctx , cancel := context.WithTimeout(context.Background(),10*time.Second)
	var project []*model.Project
	defer cancel()

	res , err := collection.Find(ctx , bson.M{})
	if err != nil {
		return nil , err
	}

	defer res.Close(ctx)

	for res.Next(ctx){
		var singleproject *model.Project
		if err = res.Decode(&singleproject); err != nil {
			log.Fatal(err)
		}
		project = append(project , singleproject)
	}

	return project, nil
}

func (db *DB) SingleOwner(ID string) (*model.Owner , error){
	collection := colHelper(db , "Owner")
	
	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var owner *model.Owner
	
	defer cancel()

	objId,_ := primitive.ObjectIDFromHex(ID)

	err := collection.FindOne(ctx , bson.M{"_id" : objId}).Decode(&owner)

	return owner,err
}

func (db *DB) SingleProject(ID string) (*model.Project , error) {
	collection := colHelper(db, "project")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var project *model.Project
    defer cancel()

	objId, _ := primitive.ObjectIDFromHex(ID)

	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&project)

	return project , err
}