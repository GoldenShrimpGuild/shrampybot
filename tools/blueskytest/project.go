package main

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	projectFilename = "../../project.yml"
	functionName    = "api"
)

type Project struct {
	Parameters  map[string]any    `yaml:"parameters"`
	Environment map[string]string `yaml:"environment"`
	Packages    []*Package        `yaml:"packages"`
}

type Package struct {
	Name        string            `yaml:"name"`
	Shared      bool              `yaml:"shared"`
	Environment map[string]string `yaml:"environment"`
	Parameters  map[string]any    `yaml:"parameters"`
	Annotations map[string]any    `yaml:"annotations"`
	Functions   []*Function       `yaml:"functions"`
	Actions     []*Function       `yaml:"actions"`
}

type Function struct {
	Name        string            `yaml:"name"`
	Runtime     string            `yaml:"runtime"`
	Web         string            `yaml:"web"`
	WebSecure   bool              `yaml:"webSecure"`
	Parameters  map[string]any    `yaml:"parameters"`
	Annotations map[string]any    `yaml:"annotations"`
	Environment map[string]string `yaml:"environment"`
	Limits      map[string]any    `yaml:"limits"`
}

func readProject(projectPath *string) (*Project, error) {
	buf, err := os.ReadFile(*projectPath)
	if err != nil {
		return nil, err
	}

	c := &Project{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", *projectPath, err)
	}

	return c, err
}

type GetEnvError error

// Merge relevant environment variables for easy lookup.
func (p *Project) getAllEnv(packageName *string) (map[string]string, error) {
	env := p.Environment
	err := errors.New("could not find package/api")

	for _, pk := range p.Packages {
		if pk.Name == *packageName {
			for k, v := range pk.Environment {
				env[k] = v
			}
			for _, f := range pk.Actions {
				err = nil
				if f.Name == functionName {
					for k, v := range f.Environment {
						env[k] = v
					}
				}
			}
			break
		}
	}

	return env, err
}
