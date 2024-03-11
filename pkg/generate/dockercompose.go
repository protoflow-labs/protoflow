package generate

import (
	"github.com/docker/cli/cli/compose/types"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"gopkg.in/yaml.v3"
	"os"
)

type DeploymentGenerator interface {
	Generate(project *project.Project) error
}

type DockerComposeGenerate struct {
	code bucket.Bucket
}

var _ DeploymentGenerator = &DockerComposeGenerate{}

func NewDockerComposeGenerate(c bucket.Bucket) *DockerComposeGenerate {
	return &DockerComposeGenerate{
		code: c,
	}
}

func (g *DockerComposeGenerate) Generate(project *project.Project) error {
	var services types.Services

	//for _, r := range project.Resources {
	//	switch r.Type.(type) {
	//	case *gen.Resource_LanguageService:
	//	}
	//}

	dataVolume := types.ServiceVolumeConfig{
		Source: "./data",
		Target: "/data",
	}

	// Add protoflow service
	services = append(services, types.ServiceConfig{
		Name:    "protoflow",
		Image:   "protoflow-labs/protoflow",
		Command: []string{"protoflow", "run"},
		Ports: []types.ServicePortConfig{
			{
				Target:    8000,
				Published: 8000,
			},
		},
		Volumes: []types.ServiceVolumeConfig{
			dataVolume,
		},
	})

	// Add postgres service
	services = append(services, types.ServiceConfig{
		Name:    "postgres",
		Image:   "postgres",
		Command: []string{"postgres"},
		Ports: []types.ServicePortConfig{
			{
				Target:    5432,
				Published: 5432,
			},
		},
		Volumes: []types.ServiceVolumeConfig{
			{
				Source: "./data/postgres",
				Target: "/var/lib/postgresql/data",
			},
		},
	})

	// Add javascript language service
	services = append(services, types.ServiceConfig{
		Name: "nodejs",
		Build: types.BuildConfig{
			Context: "./nodejs",
		},
		Command: []string{"node", "index.js"},
		Ports: []types.ServicePortConfig{
			{
				Target:    8000,
				Published: 8000,
			},
		},
	})

	out, err := yaml.Marshal(services)
	if err != nil {
		return errors.Wrapf(err, "error marshalling services")
	}

	dockerCompostPath, err := g.code.GetFile("docker-compose.yaml")
	if err != nil {
		return errors.Wrapf(err, "error getting docker-compose.yaml")
	}

	return os.WriteFile(dockerCompostPath, out, 0644)
}
