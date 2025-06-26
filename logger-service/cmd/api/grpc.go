package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/abhishek0948/goservices/logger-service/data"
	"github.com/abhishek0948/goservices/logger-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry();

	logEntry := data.LogEntry {
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry);
	if err!=nil {
		res:= &logs.LogResponse{Result: "Failed"}
		return res,err
	}

	res := &logs.LogResponse{
		Result: "logged!",
	}

	return res,nil
}

func (app *Config) gRPCListen() {
	lis,err := net.Listen("tcp",fmt.Sprintf(":%s",gRpcPort))
	if err!=nil {
		log.Fatal("Failed to listen grpc:",err);
		return ;
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s,&LogServer{Models: app.Models});

	log.Printf("Grpc stated on post:%s",gRpcPort);

	if err:= s.Serve(lis) ; err!=nil {
		log.Fatal("Failed to listen grpc:",err);
		return 
	}
}