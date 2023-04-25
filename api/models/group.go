package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v2"
)

type Group struct {
	ID           int64     `json:"id" yaml:"-"`
	Name         string    `json:"name" yaml:"name"`
	NumTargets   int64     `json:"num_targets,omitempty" yaml:"-"`
	Targets      []*Target `json:"targets,omitempty" yaml:"targets,omitempty"`
	TargetsFile  string    `json:"-" yaml:"targets-file,omitempty"`
	ModifiedDate time.Time `json:"modified_date" yaml:"-"`
	varsReplaced bool
}

func (g *Group) ToJSON() (string, error) {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (g *Group) replaceVars(vars map[string]string) error {
	if g.varsReplaced {
		return nil
	}

	name, err := templateReplace(g.Name, vars)
	if err != nil {
		return err
	}
	g.Name = name

	for _, target := range g.Targets {
		err = target.replaceVars(vars)
		if err != nil {
			return err
		}
	}

	g.varsReplaced = true
	return nil
}

func GroupFromFile(file string, vars map[string]string) (*Group, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var group Group
	err = yaml.Unmarshal(bytes, &group)
	if err != nil {
		return nil, err
	}

	parentDir := filepath.Dir(file)
	if len(group.Targets) == 0 && group.TargetsFile != "" {
		targetsPath := filepath.Join(parentDir, group.TargetsFile)

		var targets []*Target
		if filepath.Ext(targetsPath) == ".csv" {
			targetsFile, err := os.Open(targetsPath)
			if err != nil {
				return nil, err
			}
			defer targetsFile.Close()

			err = gocsv.Unmarshal(targetsFile, &targets)
			if err != nil {
				return nil, err
			}

			if len(targets) == 0 {
				targetsFile.Seek(0, 0)
				err = gocsv.UnmarshalWithoutHeaders(targetsFile, &targets)
				if err != nil {
					return nil, err
				}
			}
		} else {
			bytes, err := os.ReadFile(targetsPath)
			if err != nil {
				return nil, err
			}

			err = yaml.Unmarshal(bytes, &targets)
			if err != nil {
				return nil, err
			}
		}

		group.Targets = targets
	}

	err = group.replaceVars(vars)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

type Target struct {
	FirstName string `json:"first_name" yaml:"first-name" csv:"First Name"`
	LastName  string `json:"last_name" yaml:"last-name" csv:"Last Name"`
	Email     string `json:"email" yaml:"email" csv:"Email"`
	Position  string `json:"position" yaml:"position" csv:"Position"`
}

func (t *Target) replaceVars(vars map[string]string) error {
	email, err := templateReplace(t.Email, vars)
	if err != nil {
		return err
	}
	t.Email = email

	firstName, err := templateReplace(t.FirstName, vars)
	if err != nil {
		return err
	}
	t.FirstName = firstName

	lastName, err := templateReplace(t.LastName, vars)
	if err != nil {
		return err
	}
	t.LastName = lastName

	position, err := templateReplace(t.Position, vars)
	if err != nil {
		return err
	}
	t.Position = position

	return nil
}

type GroupsSummary struct {
	Total  int      `json:"total"`
	Groups []*Group `json:"groups"`
}
