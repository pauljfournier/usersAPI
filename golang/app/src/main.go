package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"test/api"
	"test/utils"
)

func main() {
	fmt.Println("Starting API")
	//connect to the database
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)

	defer client.Disconnect(ctx)
	awsomeDb := client.Database("awsomeDb")

	dbConnection := utils.DbConnection{
		Client:   client,
		Database: awsomeDb,
		Ctx:      ctx,
	}

	//init the server
	server, err := api.NewServer(dbConnection, false)
	if err != nil {
		log.Fatal(err)
	}
	//start the server
	server.Start()
}
