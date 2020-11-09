package main

import(
	"context"
	"fmt"
	"net/http"
	"log"
	"time"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"

)

type Article struct{
	ID           primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title        string               `json:"title,omitempty" bson:"title,omitempty"`
	Subtitle     string               `json:"subtitle,omitempty" bson:"subtitle,omitempty"`
	Content      string               `json:"content,omitempty" bson:"content,omitempty"`
}

var client *mongo.Client

func CreateArticleEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type", "application/json")
	var article Article
	json.NewDecoder(request.Body).Decode(&article)
	collection := client.Database("articleData").Collection("contentOfArticle")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, article)
	json.NewEncoder(response).Encode(result)
}

func GetArticleDataEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type", "application/json")
	var contentOfArticle []Article
	collection := client.Database("articleData").Collection("contentOfArticle")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err :=collection.Find(ctx, bson.M{})
	if err != nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(` {"message": "` + err.Error() + `"}`))
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var article Article
		cursor.Decode(&article)
		contentOfArticle = append(contentOfArticle, article)
	}
	if err:= cursor.Err(); err!=nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return	
	}
	
	json.NewEncoder(response).Encode(contentOfArticle)
}

func GetArticleEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _:= primitive.ObjectIDFromHex(params["id"])
	var article Article
	collection := client.Database("articleData").Collection("contentOfArticle")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Article{ID: id}).Decode(&article)
	if err != nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(article)
}


func main(){
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()

	router.HandleFunc("/articles", CreateArticleEndpoint).Methods("POST")
	router.HandleFunc("/articles", GetArticleDataEndpoint).Methods("GET")
	router.HandleFunc("/articles/{id}", GetArticleEndpoint).Methods("GET")
	
	log.Fatal(http.ListenAndServe(":12345", router))
	
}











