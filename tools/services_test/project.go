package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	projectDevFilename  = "../../deploy-dev/template.yml"
	projectProdFilename = "../../deploy-prod/template.yml"
)

type Template struct {
	Resources *Resources `yaml:"Resources"`
}

type Resources struct {
	ShrampyBotDev  *Resource `yaml:"shrampyBotDev"`
	ShrampyBotProd *Resource `yaml:"shrampyBotProd"`
}

type Environment struct {
	Variables map[string]string `yaml:"Variables"`
}

type Resource struct {
	Properties *Properties `yaml:"Properties"`
}

type Properties struct {
	Environment *Environment `yaml:"Environment"`
}

func readProject(deploymentEnv *string) (*Template, error) {
	var envPath string
	switch *deploymentEnv {
	case "shrampybot-dev":
		envPath = projectDevFilename
	case "shrampybot-prod":
		envPath = projectProdFilename
	}
	buf, err := os.ReadFile(envPath)
	if err != nil {
		return nil, err
	}

	c := &Template{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", envPath, err)
	}

	return c, err
}

type GetEnvError error

// Merge relevant environment variables for easy lookup.
func (p *Template) getAllEnv(deploymentEnv *string) (map[string]string, error) {
	env := map[string]string{}
	var resource *Resource

	switch *deploymentEnv {
	case "shrampybot-dev":
		resource = p.Resources.ShrampyBotDev
	case "shrampybot-prod":
		resource = p.Resources.ShrampyBotProd
	}

	for k, v := range resource.Properties.Environment.Variables {
		env[k] = v
	}

	return env, nil
}
