package models

import (
	"encoding/json"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Campaign struct {
	ID            int64           `json:"id" yaml:"-"`
	Name          string          `json:"name" yaml:"name"`
	CreatedDate   time.Time       `json:"created_date,omitempty" yaml:"-"`
	LaunchDate    time.Time       `json:"launch_date" yaml:"launch-date"`
	SendByDate    time.Time       `json:"send_by_date" yaml:"send-by-date"`
	CompletedDate time.Time       `json:"completed_date,omitempty" yaml:"-"`
	Template      *Template       `json:"template" yaml:"template"`
	Page          *Page           `json:"page" yaml:"page"`
	Status        string          `json:"status,omitempty" yaml:"status"`
	Results       []*Result       `json:"results,omitempty" yaml:"results"`
	Groups        []*Group        `json:"groups" yaml:"groups"`
	Timeline      []*Event        `json:"timeline,omitempty" yaml:"timeline"`
	SMTP          *SendingProfile `json:"smtp" yaml:"smtp"`
	URL           string          `json:"url" yaml:"url"`
	varsReplaced  bool
}

func (c *Campaign) ToJSON() (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Campaign) replaceVars(vars map[string]string) error {
	if c.varsReplaced {
		return nil
	}

	name, err := templateReplace(c.Name, vars)
	if err != nil {
		return err
	}
	c.Name = name

	url, err := templateReplace(c.URL, vars)
	if err != nil {
		return err
	}
	c.URL = url

	if c.Template != nil {
		name, err = templateReplace(c.Template.Name, vars)
		if err != nil {
			return err
		}
		c.Template.Name = name
	}

	if c.Page != nil {
		name, err = templateReplace(c.Page.Name, vars)
		if err != nil {
			return err
		}
		c.Page.Name = name
	}

	if c.SMTP != nil {
		name, err = templateReplace(c.SMTP.Name, vars)
		if err != nil {
			return err
		}
		c.SMTP.Name = name
	}

	for _, group := range c.Groups {
		name, err = templateReplace(group.Name, vars)
		if err != nil {
			return err
		}
		group.Name = name
	}

	c.varsReplaced = true
	return nil
}

func CampaignFromFile(file string, vars map[string]string) (*Campaign, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var campaign Campaign
	err = yaml.Unmarshal(bytes, &campaign)
	if err != nil {
		return nil, err
	}

	err = campaign.replaceVars(vars)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

type Result struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Position  string  `json:"position"`
	Status    string  `json:"status"`
	IP        string  `json:"ip"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	SendDate  string  `json:"send_date"`
	Reported  bool    `json:"reported"`
}

type Event struct {
	Email   string `json:"email"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Details string `json:"details"`
}
