package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/abhishek0948/goservices/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient,err := connectToMongo();
	if err!= nil {
		log.Fatal("Failed to Connect to Mongo....")
	}
	log.Println("Connected to MongoDB...")
	client = mongoClient

	ctx,cancel := context.WithTimeout(context.Background(),15*time.Second);
	defer cancel();

	defer func() {
		if err = client.Disconnect(ctx) ; err != nil {
			panic(err);
		}
	} ()

	app := Config {
		Models : data.New(client),
	}

	// Register RPC server
	err = rpc.Register(new(RPCServer));
	go app.rpcListen()
	
	// Grpc connection
	go app.gRPCListen()
	
	srv := &http.Server {
		Addr : fmt.Sprintf(":%s",webPort),
		Handler : app.routes(),
	}

	fmt.Println("Server Started");
	err = srv.ListenAndServe();
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting rpc on port:",rpcPort)
	listen,err := net.Listen("tcp",fmt.Sprintf("0.0.0.0:%s",rpcPort));
	if err!= nil {
		return err;
	}
	defer listen.Close();

	for {
		rpcConn,err := listen.Accept();
		if err!= nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client,error) {
	clientOptions := options.Client().ApplyURI(mongoURL);
	clientOptions.SetAuth(options.Credential{
		Username : "admin",
		Password : "password",
	})

	c,err := mongo.Connect(context.TODO(),clientOptions);
	if err != nil {
		log.Println("Error Connecting....",err);
		return nil,err;
	}

	return c,nil;
}