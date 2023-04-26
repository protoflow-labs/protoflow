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
				Url: "mem://",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Init()
	if err != nil {
		t.Fatal(err)
	}

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

	b, cleanup, err := resource.WithPath("test")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	err = b.WriteAll(context.Background(), "test", []byte("test"), nil)
	if err != nil {
		t.Fatal(err)
	}
	data, err := b.ReadAll(context.Background(), "test")
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
				Url: "mem://collection",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Init()
	if err != nil {
		t.Fatal(err)
	}

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

	d, _, err := resource.WithCollection("Name")
	if err != nil {
		t.Fatal(err)
	}

	type Player struct {
		Name  string
		Score int
	}

	createPlayer := Player{
		Name:  "chris",
		Score: 100,
	}

	err = d.Create(context.Background(), &createPlayer)
	if err != nil {
		t.Fatal(err)
	}

	getPlayer := Player{
		Name: "chris",
	}
	err = d.Get(context.Background(), &getPlayer)
	if err != nil {
		t.Fatal(err)
	}
	if createPlayer.Score != getPlayer.Score {
		t.Fatalf("expected %d, got %d", createPlayer.Score, getPlayer.Score)
	}
}
