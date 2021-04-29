package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sepehrhakimi90/GOPostBackend/entity"
)

// TODO: DB Config validation

const (
	postCollection = "posts"
)

var (
	dbName = ""
	dbUsername = ""
	dbPassword = ""
	dbURL = ""
)


type mongoPostRepo struct {
	mongoClient *mongo.Client
}

func (m mongoPostRepo) Save(post *entity.Post) (*entity.Post, error) {
	db := m.mongoClient.Database(dbName)
	collection := db.Collection(postCollection)
	result, err := collection.InsertOne(getContext(), post)
	if err != nil {
		log.Println("Error in insert data", err)
		return nil, err
	}
	id, _:= result.InsertedID.(primitive.ObjectID)
	post.Id = id.Hex()

	return post, nil
}

func (m mongoPostRepo) FindAll() ([]entity.Post, error) {
	posts := make([]entity.Post,0)
	db := m.mongoClient.Database(dbName)
	collection := db.Collection(postCollection)

	cur, err := collection.Find(getContext(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var post entity.Post
		err := cur.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (m mongoPostRepo) FindById(id string) (*entity.Post, error) {
	post := entity.Post{}
	db := m.mongoClient.Database(dbName)
	collection := db.Collection(postCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Mongo-Repository", err)
		return nil, err
	}
	res := collection.FindOne(getContext(), bson.M{"_id": objectID})
	err = res.Decode(&post)
	if err != nil {
		log.Println("Mongo-Repository", err)
		return nil, err
	}

	return &post, nil
}

//NewMongoPostRepo

func NewMongoPostRepo() (PostRepository, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return &mongoPostRepo{mongoClient: client}, nil
}

func getClient() (*mongo.Client, error) {
	mongoServerIP := os.Getenv("MONGO_IP")
	mongoServerPort := os.Getenv("MONGO_PORT")
	dbURL = fmt.Sprintf("mongodb://%s:%s", mongoServerIP, mongoServerPort)
	dbUsername = os.Getenv("MONGO_USER")
	dbPassword = os.Getenv("MONGO_PASSWORD")
	dbName = os.Getenv("MONGO_DB_NAME")

	clientOption := options.Client().ApplyURI(dbURL).SetAuth(options.Credential{Username: dbUsername, Password: dbPassword, AuthSource: dbName})
	mongoClient, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		return nil, err
	}

	if err = mongoClient.Ping(getContext(), nil); err != nil {
		return nil, err
	}

	return mongoClient, nil

}

func getContext() context.Context{
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}
