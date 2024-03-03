package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/SunnyChugh99/banking_management_golang/doc/statik"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/gapi"
	"github.com/SunnyChugh99/banking_management_golang/pb"
	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	store := db.NewStore(conn)
	go runHTTPGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot create listen", err)
	}

	log.Printf("Start grpc server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start grpc server", err)
	}
}


func runHTTPGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	// grpcServer := grpc.NewServer()
	// pb.RegisterSimpleBankServer(grpcServer, server)
	// reflection.Register(grpcServer)

	gprcMutex := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, gprcMutex, server)
	if err != nil {
		log.Fatal("Cannot create handler", err)
	}

	Mux := http.NewServeMux()

	Mux.Handle("/", gprcMutex)

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("Cannot create statik fs", err)
	}

	swaggerHandler := http.StripPrefix("/swagger", http.FileServer(statikFs))
	Mux.Handle("/swagger/", swaggerHandler)

	
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot create listen", err)
	}

	log.Printf("Start HTTP server at %s", listener.Addr().String())

	err = http.Serve(listener, Mux)
	if err != nil {
		log.Fatal("Cannot start HTTP server", err)
	}
}
