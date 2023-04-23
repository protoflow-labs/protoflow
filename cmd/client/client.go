package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
)

func main() {
	client := genconnect.NewProjectServiceClient(
		http.DefaultClient,
		"http://localhost:8080/",
	)

	req := connect.NewRequest(&gen.GetProjectsRequest{})
	req.Header().Set("Some-Header", "hello from connect")
	res, err := client.GetProjects(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
	log.Println(res.Header().Get("Some-Other-Header"))
}
