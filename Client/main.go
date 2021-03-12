package main

import (
	"encoding/json"
	"log"
	"net/http"
	"context"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address = "grpc-server:4000"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func newElement(w http.ResponseWriter, r *http.Request) {
	// Adding headers
	w.Header().Set("Content-Type", "application/json")

	// Parsing body
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	failOnError(err, "Parsing JSON")
	body["way"] = "GRPC"
	data, err := json.Marshal(body)

	newData := string(data)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err, "GRPC Connection")
	defer conn.Close()

	cli := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	re, err := cli.SayHello(ctx, &pb.HelloRequest{Name: string(data)})
	if err != nil {
		failOnError(err, "Error al enviar el mensaje")
	}
	log.Print("Enviado")
	log.Printf("respuesta : %s", re.GetMessage())

	// Setting status and send response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(newData))
}

func handleRequests() {
	http.HandleFunc("/", newElement)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	handleRequests()
}