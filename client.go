package okrforjira

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	jsontime "github.com/liamylian/jsontime/v2/v2"
	"golang.org/x/exp/slices"
)

// Client is a HTTP client to get data from the OKR for Jira app.
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient creates a new OKR for Jira client.
// If a nil httpClient is provided, a new http.Client will be used.
func NewClient(httpClient *http.Client, token string) *Client {
	if httpClient == nil {
		// Avoid to use the default transport.
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxConnsPerHost = 100
		t.MaxIdleConnsPerHost = 100
		httpClient = &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		}
	}

	return &Client{
		httpClient: httpClient,
		token:      token,
	}
}

type Response struct {
	OKRs       []OKR       `json:"okrs"`
	KeyResults []KeyResult `json:"krs"`
	Teams      []Team      `json:"teams"`
	Periods    []Period    `json:"periods"`
	Labels     []Label     `json:"labels"`
}

type OKR struct {
	ID                     string    `json:"id"`
	Key                    string    `json:"key"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	ParentObjectiveID      string    `json:"parentObjectiveId"`
	OwnerAccountID         string    `json:"ownerAccountId"`
	CollaboratorAccountIDs []string  `json:"collaboratorAccountIds"`
	PercentDone            float64   `json:"percentDone"`
	Created                time.Time `json:"created" time_format:"okr4j_format"`
	StartDate              time.Time `json:"startDate" time_format:"okr4j_format"`
	Deadline               time.Time `json:"deadline" time_format:"okr4j_format"`
	LabelIDs               []string  `json:"labelIds"`
	TeamIDs                []string  `json:"teamIds"`
	KRIDs                  []string  `json:"krIds"`
	ChildObjectiveIDs      []string  `json:"childObjectiveIds"`
	LatestUpdate           Update    `json:"latestUpdate"`
	PeriodAliasID          string    `json:"periodAliasId"`
	Weight                 float64   `json:"weight"`
}

type KeyResult struct {
	ID                        string             `json:"id"`
	Key                       string             `json:"key"`
	Name                      string             `json:"name"`
	Description               string             `json:"description"`
	ParentObjectiveID         string             `json:"parentObjectiveId"`
	IssueIDs                  []string           `json:"issueIds"`
	OwnerAccountID            string             `json:"ownerAccountId"`
	CollaboratorAccountIds    []string           `json:"collaboratorAccountIds"`
	PercentDone               float64            `json:"percentDone"`
	Created                   time.Time          `json:"created" time_format:"okr4j_format"`
	StartDate                 time.Time          `json:"startDate" time_format:"okr4j_format"`
	Deadline                  time.Time          `json:"deadline" time_format:"okr4j_format"`
	LabelIDs                  []string           `json:"labelIds"`
	TeamIDs                   []string           `json:"teamIds"`
	PeriodAliasID             string             `json:"periodAliasId"`
	LatestUpdate              Update             `json:"latestUpdate"`
	Unit                      Unit               `json:"unit"`
	CurrentProgressDefinition ProgressDefinition `json:"currentProgressDefinition"`
	Weight                    float64            `json:"weight"`
}

type Update struct {
	EntityID    string    `json:"entityId"`
	Status      string    `json:"status"`
	Created     time.Time `json:"created" time_format:"okr4j_format"`
	Value       float64   `json:"value"`
	Description string    `json:"description"`
}

type ProgressDefinition struct {
	Type         string  `json:"type"`
	StartValue   float64 `json:"startValue"`
	DesiredValue float64 `json:"desiredValue"`
	JQL          string  `json:"jql"`
}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Period struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate" time_format:"okr4j_format"`
	Deadline  time.Time `json:"deadline" time_format:"okr4j_format"`
}

type Label struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Unit struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// ObjectivesByDate returns a list of objectives which have start date or/and due date inside specified date range.
func (c *Client) ObjectivesByDate(ctx context.Context, startDate, deadline time.Time, expand []string) (Response, error) {
	startDateEpochMilli := startDate.UnixMilli()
	deadlineEpochMilli := deadline.UnixMilli()
	if err := c.checkObject(expand); err != nil {
		return Response{}, fmt.Errorf("failed to get the obectives by date: %w", err)
	}
	url := fmt.Sprintf(objectivesByDateURL, startDateEpochMilli, deadlineEpochMilli, strings.Join(expand, ","))
	response, err := c.executeGetQuery(url)
	if err != nil {
		return Response{}, fmt.Errorf("failed to get the obectives by date: %w", err)
	}
	return response, nil
}

// ObjectivesByIDs returns a list of objectives with specified ids.
func (c *Client) ObjectivesByIDs(ctx context.Context, objectiveIDs, expand []string) (Response, error) {
	if err := c.checkObject(expand); err != nil {
		return Response{}, fmt.Errorf("failed to get the obectives by ids: %w", err)
	}
	url := fmt.Sprintf(objectivesByIDsURL, strings.Join(objectiveIDs, ","), strings.Join(expand, ","))
	response, err := c.executeGetQuery(url)
	if err != nil {
		return Response{}, fmt.Errorf("failed to get the obectives by ids: %w", err)
	}
	return response, nil
}

// KeyResultsByDate returns a list of key results which have start date or/and due date inside specified date range.
func (c *Client) KeyResultsByDate(ctx context.Context, startDate, deadline time.Time, expand []string) (Response, error) {
	startDateEpochMilli := startDate.UnixMilli()
	deadlineEpochMilli := deadline.UnixMilli()
	if err := c.checkObject(expand); err != nil {
		return Response{}, fmt.Errorf("failed to get the obectives by date: %w", err)
	}
	url := fmt.Sprintf(keyResultsByDateURL, startDateEpochMilli, deadlineEpochMilli, strings.Join(expand, ","))
	response, err := c.executeGetQuery(url)
	if err != nil {
		return Response{}, fmt.Errorf("failed to get the obectives by date: %w", err)
	}
	return response, nil
}

// KeyResultsByIDs returns a list of key results with specified ids.
func (c *Client) KeyResultsByIDs(ctx context.Context, keyResultIDs, expand []string) (Response, error) {
	if err := c.checkObject(expand); err != nil {
		return Response{}, fmt.Errorf("failed to get the key results by ids: %w", err)
	}
	url := fmt.Sprintf(keyREsultsByIDsURL, strings.Join(keyResultIDs, ","), strings.Join(expand, ","))
	response, err := c.executeGetQuery(url)
	if err != nil {
		return Response{}, fmt.Errorf("failed to get the key results by ids: %w", err)
	}
	return response, nil
}

type objectiveUpdateRequest struct {
	ObjectiveID string `json:"objectiveId"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// UpdateObjective update the provided objective.
func (c *Client) UpdateObjective(ctx context.Context, objectiveID, status, description string) (Update, error) {
	request := objectiveUpdateRequest{
		ObjectiveID: objectiveID,
		Status:      status,
		Description: description,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return Update{}, fmt.Errorf("failed to update the objective: %w", err)
	}
	var response Update
	if err := c.executePostQuery(objectiveUpdateURL, string(data), &response); err != nil {
		return Update{}, fmt.Errorf("failed to update the objective: %w", err)
	}
	return response, nil
}

type keyResultUpdateRequest struct {
	KeyResultID string  `json:"keyResultId"`
	Status      string  `json:"status"`
	NewValue    float64 `json:"newValue"`
	Description string  `json:"description"`
}

// UpdateKeyResult update the provided key result.
func (c *Client) UpdateKeyResult(ctx context.Context, keyResultID, status string, newValue float64, description string) (Update, error) {
	request := keyResultUpdateRequest{
		KeyResultID: keyResultID,
		Status:      status,
		NewValue:    newValue,
		Description: description,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return Update{}, fmt.Errorf("failed to update the key result: %w", err)
	}
	var response Update
	if err := c.executePostQuery(keyResultUpdateURL, string(data), &response); err != nil {
		return Update{}, fmt.Errorf("failed to update the key result: %w", err)
	}
	return response, nil
}

const (
	host                = "https://okr-for-jira-prod.herokuapp.com"
	objectivesByDateURL = host + "/api/v2/api-export/objectives/byDate?startDateEpochMilli=%d&deadlineEpochMilli=%d&expand=%s"
	objectivesByIDsURL  = host + "/api/v2/api-export/objectives/byIds?objectiveIds=%s&expand=%s"
	keyResultsByDateURL = host + "/api/v2/api-export/keyResults/byDate?startDateEpochMilli=%d&deadlineEpochMilli=%d&expand=%s"
	keyREsultsByIDsURL  = host + "/api/v2/api-export/keyResults/byIds?keyResultIds=%s&expand=%s"
	objectiveUpdateURL  = host + "/api/v2/api-update/objectives"
	keyResultUpdateURL  = host + "/api/v2/api-update/keyResults"
	dateFormat          = "2006-01-02T15:04:05-0700"
)

var validObjectTypes = []string{"OBJECTIVES", "KEY_RESULTS", "TEAMS", "PERIODS", "LABELS"}

func (c Client) executeGetQuery(url string) (Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("API-Token", c.token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r, err := c.httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return Response{}, err
		}
		return Response{}, fmt.Errorf("error status %d: %s", r.StatusCode, string(content))
	}
	var response Response

	json := jsontime.ConfigWithCustomTimeFormat
	jsontime.AddTimeFormatAlias("okr4j_format", dateFormat)
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return Response{}, err
	}
	return response, nil
}

func (c Client) executePostQuery(url string, body string, response interface{}) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("API-Token", c.token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 && r.StatusCode != 201 {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("error status %d: %s", r.StatusCode, string(content))
	}

	json := jsontime.ConfigWithCustomTimeFormat
	jsontime.AddTimeFormatAlias("okr4j_format", dateFormat)
	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		return err
	}
	return nil
}

func (c *Client) checkObject(expand []string) error {
	for _, object := range expand {
		if !slices.Contains(validObjectTypes, object) {
			return fmt.Errorf("invalid object %s", object)
		}
	}
	return nil
}
