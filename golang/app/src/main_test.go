package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"test/api"
	"test/utils"
	"testing"
)

var dbConnection utils.DbConnection
var server *api.Server

func TestMain(m *testing.M) {
	fmt.Println("Starting API for TEST")
	//connect to the database
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)

	defer client.Disconnect(ctx)
	testDb := client.Database("testDb")

	dbConnection = utils.DbConnection{
		Client:   client,
		Database: testDb,
		Ctx:      ctx,
	}

	//init the server
	server, err = api.NewServer(dbConnection, true)
	if err != nil {
		log.Fatal(err)
	}
	//start the server
	server.Start()

	//run the tests
	exitVal := m.Run()

	//DROP the db and stop the server to clean
	err = testDb.Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = server.Close()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

func TestPingTrueDb(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/ping")
	if err != nil {
		t.Errorf("Ping test for main db get an err %v", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Ping test for main db get a status code %v", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Ping test for main db get an err %v trying to parse the body", err.Error())
	}
	if string(b) != "pong" {
		t.Errorf("Ping test for main db get %v expected pong", string(b))
	}
}

//test ping
//test connect to

////test insert

//normal
//missing field
//non uniq email
//non uniq nickname
//non valid email
//with created_at and id

//from store
//from ressource
//from request

////test update
//normal
//non existing
//missing field
//non uniq email
//non uniq nickname
//with created_at

//from store
//from ressource
//from request

////test delete
//existing
//non existing

//from store
//from ressource
//from request

////test list
//page and page size normal
//page and page size partial page
//page and page size no result
//test text in multiple field
//page test for each parameter and 2 results, but not a third one
//test created before/after and both
//test updated before/after and both
//mega query with one result

//from store
//from ressource
//from request
