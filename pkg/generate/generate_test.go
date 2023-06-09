package generate

import (
	"testing"

	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	var err error
	project := &gen.Project{
		Name: "test",
		Graph: &gen.Graph{
			Nodes: []*gen.Node{
				{
					Name: "SayHello",
					Config: &gen.Node_Function{
						Function: &gen.Function{
							Runtime: "node",
						},
					},
				},
			},
		},
	}

	cache, err := bucket.NewUserCache(bucket.Config{
		Name: ".protoflow_test",
	})
	assert.NoError(t, err)

	generator := NewGenerate(cache)
	assert.NoError(t, err)

	err = generator.Generate(project)
	assert.NoError(t, err)
}
