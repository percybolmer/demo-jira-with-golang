package jira

import (
	"os"

	gojira "gopkg.in/andygrunwald/go-jira.v1"
)

// JIRA is a wrapper to gojira client that makes it easier to use and match the usecases we have
type JIRA struct {
	client *gojira.Client
}

func New() (*JIRA, error) {
	// Create a BasicAuth Transport object
	tp := gojira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}
	// Create a new Jira Client
	client, err := gojira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		return nil, err
	}

	return &JIRA{
		client: client,
	}, nil
}

// TransitionIssue will move a issue into the new transition
func (j *JIRA) TransitionIssue(issue gojira.Issue, transition gojira.Transition) error {
	_, err := j.client.Issue.DoTransition(issue.ID, transition.ID)
	return err
}

// GetIssueTransition will grab the available transitions for a issue
func (j *JIRA) GetIssueTransition(issue gojira.Issue, status string) (gojira.Transition, error) {
	transitions, _, err := j.client.Issue.GetTransitions(issue.Key)
	if err != nil {
		return gojira.Transition{}, err
	}
	for _, t := range transitions {
		if t.Name == status {
			return t, nil
		}
	}
	return gojira.Transition{}, nil
}

// GetIssues will query Jira API using the provided JQL string
func (j *JIRA) GetIssues(jql string) ([]gojira.Issue, error) {

	// lastIssue is the index of the last issue returned
	lastIssue := 0
	// Make a loop through amount of issues
	var result []gojira.Issue
	for {
		// Add a Search option which accepts maximum amount (1000)
		opt := &gojira.SearchOptions{
			MaxResults: 1000,      // Max amount
			StartAt:    lastIssue, // Make sure we start grabbing issues from last checkpoint
		}
		issues, resp, err := j.client.Issue.Search(jql, opt)
		if err != nil {
			return nil, err
		}
		// Grab total amount from response
		total := resp.Total
		if issues == nil {
			// init the issues array with the correct amount of length
			result = make([]gojira.Issue, 0, total)
		}

		// Append found issues to result
		result = append(result, issues...)
		// Update checkpoint index by using the response StartAt variable
		lastIssue = resp.StartAt + len(issues)
		// Check if we have reached the end of the issues
		if lastIssue >= total {
			break
		}
	}
	return result, nil
}

// GetProjects is a function that lists all projects
func (j *JIRA) GetProjects() (*gojira.ProjectList, error) {
	// use the Project domain to list all projects
	projectList, _, err := j.client.Project.GetList()
	if err != nil {
		return nil, err
	}
	return projectList, nil
}
