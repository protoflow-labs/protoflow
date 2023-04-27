package generate

import (
	"testing"

	"github.com/protoflow-labs/protoflow/gen"
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

	generator, err := NewGenerate(project)
	assert.NoError(t, err)

	err = generator.Generate()
	assert.NoError(t, err)
}
