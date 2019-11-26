package main

import (
	"context"
	"net/http"
	"log" 
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "io/ioutil"
)

//Information received from MongoDB
type ReceivedInformation struct {
	All_repos_number int
	Google_repos_number int
	Microsoft_repos_number int
	
	Google_total_lines_of_code int
	Google_contributors Contributors
	
	Microsoft_total_lines_of_code int
	Microsoft_contributors Contributors
}

//describes the information we need about contributors
type Contributors struct {	
	Employee_count int
	Non_employee_count int
	Employees_line_count int
	Non_employees_line_count int
	
	Employee_languages []Languages
	Non_employee_languages []Languages
}

//information about languages used
type Languages struct {
	Name string
	Lines_of_changes int
}

func main() {
	// get mongoDB username and password
	m_username, err := ioutil.ReadFile("src/frontEnd/username.txt") // file with just mongoDB username in it
	if err != nil {
    	log.Fatal(err)
    }
	m_password, err := ioutil.ReadFile("src/frontEnd/password.txt") // file with just mongoDB password in it
	if err != nil {
    	log.Fatal(err) 
    }
	URI := "mongodb+srv://" + string(m_username) + ":" + string(m_password) + "@sweng-blmoo.azure.mongodb.net/test?retryWrites=true&w=majority"
	
	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(URI)
	mongo_client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
	    log.Fatal(err)
	}
	// Check the connection
	err = mongo_client.Ping(context.Background(), nil)
	if err != nil {
	    log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	collection := mongo_client.Database("Software_Engineering").Collection("Cached_Data")
	
	filter := bson.D{{}}
	var data ReceivedInformation

	err = collection.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
	    log.Fatal(err)
	}
	log.Println("Data fetched")

	
	http.HandleFunc("/index/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/frontEnd/static"))))
    log.Fatal(http.ListenAndServe(":8080", nil))
    
    
}