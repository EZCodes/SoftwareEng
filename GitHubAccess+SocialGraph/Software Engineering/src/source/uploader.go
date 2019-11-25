package main 

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "log"
    "time"
    "fmt"
)

// one big struct that we upload to MongoDB
type InformationToUpload struct {
	all_repos_number int
	google_repos_number int
	microsoft_repos_number int
	
	google_total_lines_of_code int
	google_contributors ContributorsInformation
	
	microsoft_total_lines_of_code int
	microsoft_contributors ContributorsInformation
}

//describes the information we need about contributors
type ContributorsInformation struct {	
	employee_count int
	non_employee_count int
	employees_line_count int
	non_employees_line_count int
	
	employee_languages []Language
	non_employee_languages []Language
}

//information about languages used
type Language struct {
	name string
	lines_of_changes int
}

// Will upload out information to the MongoDB
// To simplify the process, first we delete an entry and then upload a new one
func uploadToMongo(data, cache_chan chan InformationToUpload , collection *mongo.Collection) {
	for {
		val, ok := <-data
		if !ok {
			break
		}
		_, err := collection.DeleteOne(context.Background(), bson.D{{}})
		if err != nil {
		    log.Fatal(err)
		}
		_, err = collection.InsertOne(context.Background(), val)
		if err != nil {
		    log.Fatal(err)
		}
		// soft push into channel (doesnt block)
		select {
			case cache_chan <- val: // push item it channel empty
		    default: // dont push if chanell is full
		}
	}
}

// Every hour update our "cache" collection
func uploadToMongoCache(data chan InformationToUpload , collection *mongo.Collection) {
	for {
		val, ok := <-data
		if !ok {
			break
		}
		deleteResult, err := collection.DeleteOne(context.Background(), bson.D{{}})
		if err != nil {
		    log.Fatal(err)
		}
		log.Println("Succesfully deleted an entry %d", deleteResult.DeletedCount)
		insertResult, err := collection.InsertOne(context.Background(), val)
		if err != nil {
		    log.Fatal(err)
		}
		log.Println("Uploaded new data Successfully %s", insertResult.InsertedID)
		time.Sleep(time.Hour)
	}
}