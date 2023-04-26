package workflow

import (
	"context"
	"github.com/protoflow-labs/protoflow/gen"
	"testing"
)

func TestBlobstoreResource_Init(t *testing.T) {
	r, err := ResourceFromProto(&gen.Resource{
		Id: "1",
		Type: &gen.Resource_Blobstore{
			Blobstore: &gen.Blobstore{
				Url: "mem://collection/test",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	cleanup, err := r.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	resources := map[string]any{
		"1": r,
	}

	// try to look up wrong resource
	_, err = getResource[DocstoreResource](resources)
	if err == nil {
		t.Fatal("expected error")
	}

	resource, err := getResource[BlobstoreResource](resources)
	if err != nil {
		t.Fatal(err)
	}
	if resource == nil {
		t.Fatal("resource is nil")
	}
	err = resource.Bucket.WriteAll(context.Background(), "test", []byte("test"), nil)
	if err != nil {
		t.Fatal(err)
	}
	data, err := resource.Bucket.ReadAll(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "test" {
		t.Fatalf("expected %s, got %s", "test", string(data))
	}
}

func TestDocstoreResource_Init(t *testing.T) {
	r, err := ResourceFromProto(&gen.Resource{
		Id: "1",
		Type: &gen.Resource_Docstore{
			Docstore: &gen.Docstore{
				Url: "mem://collection/Name",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	cleanup, err := r.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	resources := map[string]any{
		"1": r,
	}

	resource, err := getResource[DocstoreResource](resources)
	if err != nil {
		t.Fatal(err)
	}
	if resource == nil {
		t.Fatal("resource is nil")
	}

	type Player struct {
		Name  string
		Score int
	}

	createPlayer := Player{
		Name:  "chris",
		Score: 100,
	}

	err = resource.Collection.Create(context.Background(), &createPlayer)
	if err != nil {
		t.Fatal(err)
	}

	getPlayer := Player{
		Name: "chris",
	}
	err = resource.Collection.Get(context.Background(), &getPlayer)
	if err != nil {
		t.Fatal(err)
	}
	if createPlayer.Score != getPlayer.Score {
		t.Fatalf("expected %d, got %d", createPlayer.Score, getPlayer.Score)
	}
}
