package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type CreateIssueRequestParams struct {
	Fields FieldsRequestParams `json:"fields"`
}

type FieldsRequestParams struct {
	Project     ProjectRequestParams   `json:"project"`
	Assignee    AssigneeRequestParams  `json:"assignee"`
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Issuetype   IssuetypeRequestParams `json:"issuetype"`
}

type ProjectRequestParams struct {
	Key string `json:"key"`
}

type AssigneeRequestParams struct {
	Name string `json:"name"`
}

type IssuetypeRequestParams struct {
	Name string `json:"name"`
}

type CreateIssueResponseBody struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func createIssue(project string, title string, issueType string, description string) (*CreateIssueResponseBody, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	params := CreateIssueRequestParams{
		Fields: FieldsRequestParams{
			Project: ProjectRequestParams{
				Key: project,
			},
			Assignee: AssigneeRequestParams{
				Name: config.Username,
			},
			Issuetype: IssuetypeRequestParams{
				Name: issueType,
			},
			Summary:     title,
			Description: description,
		},
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(params)

	fmt.Printf("%v\n", b)

	response := &CreateIssueResponseBody{}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/rest/api/2/issue", config.JiraHost),
		b,
	)
	req.Header.Add("Content-Type", "application/json") // ;charset=utf-8
	req.SetBasicAuth(config.Username, config.Password)
	res, err := client.Do(req)

	if err == nil && res.StatusCode != 201 {
		msg := fmt.Sprintf("Wrong status code %s", res.Status)
		err = errors.New(msg)
	}
	if err != nil {
		fmt.Printf("Call failed: %s", err.Error())
		return response, err
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(response)

	fmt.Printf("Created %s/browse/%s\n", config.JiraHost, response.Key)

	return response, err
}
