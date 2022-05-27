package okrforjira_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grandper/okrforjira"
	"github.com/stretchr/testify/assert"
)

// The tests are based on the examples provided by Digital Toucan
// on the following pages:
// - https://intercom.help/okr-for-jira/en/articles/6178378-api-query-methods .
// - https://intercom.help/okr-for-jira/en/articles/6252250-api-update-methods .

func TestNewClient(t *testing.T) {
	t.Run("Create a simple client", func(t *testing.T) {
		c := okrforjira.NewClient(nil, token)
		assert.NotNil(t, c)
	})
}

const token = "foobar-token"

func TestClient_ObjectivesByDate(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "https://okr-for-jira-prod.herokuapp.com/api/v2/api-export/objectives/byDate?startDateEpochMilli=1409459200000&deadlineEpochMilli=1748647410000&expand=KEY_RESULTS,TEAMS,PERIODS,LABELS")
		assert.Equal(t, req.Header.Get("API-Token"), token)
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
	"okrs": [
		{
			"id": "5fda249d289742000406b3e4",
			"key": "O-2",
			"name": "Become more mature company",
			"description": "<p>This quarter we will be focusing on improving our performance.</p><p></p>",
			"parentObjectiveId": null,
			"ownerAccountId": "5c12ad9fd3af3b1ccfecbf55",
			"collaboratorAccountIds": [],
			"percentDone": 8.333333333333332,
			"created": "2020-12-16T15:15:41+0000",
			"startDate": "2021-01-01T00:00:00+0000",
			"deadline": "2021-03-31T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"krIds": [
				"605480b190c42b0003385170",
				"6061e921e2f4470003bc3210"
			],
			"childObjectiveIds": [
				"5fdb72c63d2cf000035ceb37",
				"60743135b347480003dc6a9c",
				"61f9367df9aa7f0e4024a6fe"
			],
			"latestUpdate": {
				"entityId": "5fda249d289742000406b3e4",
				"status": "ON_TRACK",
				"created": "2021-05-05T12:15:14+0000",
				"value": null,
				"description": ""
			},
			"periodAliasId": "602a6a2717378700039f342a",
			"weight": 0
		}
	],
	"krs": [
		{
			"id": "605480b190c42b0003385170",
			"key": "KR-8",
			"name": "new auto KR",
			"description": null,
			"parentObjectiveId": "5fda249d289742000406b3e4",
			"issueIds": [
				"10005",
				"10006",
				"10010"
			],
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [
				"5c12ad9fd3af3b1ccfecbf55"
			],
			"percentDone": 33.33333333333333,
			"created": "2021-03-19T10:45:05+0000",
			"startDate": "2021-01-01T00:00:00+0000",
			"deadline": "2021-03-31T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342a",
			"latestUpdate": {
				"entityId": "620ea538512edb00acf67ac1",
				"status": "ON_TRACK",
				"created": "2022-02-17T19:42:43+0000",
				"value": 1.0,
				"description": null
			},
			"unit": null,
			"currentProgressDefinition": {
				"type": "AUTO",
				"startValue": null,
				"desiredValue": null,
				"jql": "project= TEST"
			},
			"weight": 1
		},
		{
			"id": "6061e921e2f4470003bc3210",
			"key": "KR-9",
			"name": "different start date",
			"description": null,
			"parentObjectiveId": "5fda249d289742000406b3e4",
			"issueIds": [
				"10000"
			],
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2021-03-29T14:50:09+0000",
			"startDate": "2020-04-01T00:00:00+0000",
			"deadline": "2020-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": null,
			"latestUpdate": {
				"entityId": "61138f2be5fd454858c3e1ee",
				"status": "AT_RISK",
				"created": "2021-08-11T08:49:47+0000",
				"value": 0.0,
				"description": null
			},
			"unit": {
				"name": "USD",
				"symbol": "$"
			},
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 1
		}
	],
	"teams": [],
	"periods": [
		{
			"id": "602a6a2717378700039f342a",
			"name": "Q1 Y2021",
			"startDate": "2021-01-01T00:00:00+0000",
			"deadline": "2021-03-31T23:59:59+0000"
		}
	],
	"labels": []
}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	ctx := context.Background()

	loc, _ := time.LoadLocation("Europe/Paris")
	startDate := time.Date(2014, time.August, 31, 6, 26, 40, 0, loc)
	deadline := time.Date(2025, time.May, 31, 1, 23, 30, 0, loc)
	expand := []string{"KEY_RESULTS", "TEAMS", "PERIODS", "LABELS"}
	c := okrforjira.NewClient(client, token)
	got, err := c.ObjectivesByDate(ctx, startDate, deadline, expand)
	assert.NoError(t, err)

	want := okrforjira.Response{
		OKRs: []okrforjira.OKR{
			{
				ID:                     "5fda249d289742000406b3e4",
				Key:                    "O-2",
				Name:                   "Become more mature company",
				Description:            "<p>This quarter we will be focusing on improving our performance.</p><p></p>",
				ParentObjectiveID:      "",
				OwnerAccountID:         "5c12ad9fd3af3b1ccfecbf55",
				CollaboratorAccountIDs: []string{},
				PercentDone:            8.333333333333332,
				Created:                time.Date(2020, time.December, 16, 15, 15, 41, 0, time.UTC),
				StartDate:              time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2021, time.March, 31, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				KRIDs: []string{
					"605480b190c42b0003385170",
					"6061e921e2f4470003bc3210",
				},
				ChildObjectiveIDs: []string{
					"5fdb72c63d2cf000035ceb37",
					"60743135b347480003dc6a9c",
					"61f9367df9aa7f0e4024a6fe",
				},
				LatestUpdate: okrforjira.Update{
					EntityID:    "5fda249d289742000406b3e4",
					Status:      "ON_TRACK",
					Created:     time.Date(2021, time.May, 5, 12, 15, 14, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				PeriodAliasID: "602a6a2717378700039f342a",
				Weight:        0.0,
			},
		},
		KeyResults: []okrforjira.KeyResult{
			{
				ID:                "605480b190c42b0003385170",
				Key:               "KR-8",
				Name:              "new auto KR",
				Description:       "",
				ParentObjectiveID: "5fda249d289742000406b3e4",
				IssueIDs: []string{
					"10005",
					"10006",
					"10010",
				},
				OwnerAccountID: "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIds: []string{
					"5c12ad9fd3af3b1ccfecbf55",
				},
				PercentDone:   33.33333333333333,
				Created:       time.Date(2021, time.March, 19, 10, 45, 5, 0, time.UTC),
				StartDate:     time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
				Deadline:      time.Date(2021, time.March, 31, 23, 59, 59, 0, time.UTC),
				LabelIDs:      []string{},
				TeamIDs:       []string{},
				PeriodAliasID: "602a6a2717378700039f342a",
				LatestUpdate: okrforjira.Update{
					EntityID:    "620ea538512edb00acf67ac1",
					Status:      "ON_TRACK",
					Created:     time.Date(2022, time.February, 17, 19, 42, 43, 0, time.UTC),
					Value:       1.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "AUTO",
					StartValue:   0.0,
					DesiredValue: 0.0,
					JQL:          "project= TEST",
				},
				Weight: 1.0,
			},
			{
				ID:                "6061e921e2f4470003bc3210",
				Key:               "KR-9",
				Name:              "different start date",
				Description:       "",
				ParentObjectiveID: "5fda249d289742000406b3e4",
				IssueIDs: []string{
					"10000",
				},
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2021, time.March, 29, 14, 50, 9, 0, time.UTC),
				StartDate:              time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2020, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "",
				LatestUpdate: okrforjira.Update{
					EntityID:    "61138f2be5fd454858c3e1ee",
					Status:      "AT_RISK",
					Created:     time.Date(2021, time.August, 11, 8, 49, 47, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "USD",
					Symbol: "$",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
		},
		Teams: []okrforjira.Team{},
		Periods: []okrforjira.Period{
			{
				ID:        "602a6a2717378700039f342a",
				Name:      "Q1 Y2021",
				StartDate: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
				Deadline:  time.Date(2021, time.March, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		Labels: []okrforjira.Label{},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("unexpected result:\n%s", cmp.Diff(got, want))
	}
}

func TestClient_ObjectivesByIDs(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "https://okr-for-jira-prod.herokuapp.com/api/v2/api-export/objectives/byIds?objectiveIds=5fda249d289742000406b3e4,5fdb72c63d2cf000035ceb37&expand=OBJECTIVES,KEY_RESULTS,TEAMS,PERIODS,LABELS")
		assert.Equal(t, req.Header.Get("API-Token"), token)
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
	"okrs": [
		{
			"id": "5fda249d289742000406b3e4",
			"key": "O-2",
			"name": "Become more mature company",
			"description": "<p>This quarter we will be focusing on improving our performance.</p><p></p>",
			"parentObjectiveId": null,
			"ownerAccountId": "5c12ad9fd3af3b1ccfecbf55",
			"collaboratorAccountIds": [],
			"percentDone": 8.333333333333332,
			"created": "2020-12-16T15:15:41+0000",
			"startDate": "2021-01-01T00:00:00+0000",
			"deadline": "2021-03-31T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"krIds": [
				"605480b190c42b0003385170",
				"6061e921e2f4470003bc3210"
			],
			"childObjectiveIds": [
				"5fdb72c63d2cf000035ceb37",
				"60743135b347480003dc6a9c",
				"61f9367df9aa7f0e4024a6fe"
			],
			"latestUpdate": {
				"entityId": "5fda249d289742000406b3e4",
				"status": "ON_TRACK",
				"created": "2021-05-05T12:15:14+0000",
				"value": null,
				"description": ""
			},
			"periodAliasId": "602a6a2717378700039f342a",
			"weight": 0
		},
		{
			"id": "5fdb72c63d2cf000035ceb37",
			"key": "O-3",
			"name": "45",
			"description": null,
			"parentObjectiveId": "5fda249d289742000406b3e4",
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2020-12-17T15:01:26+0000",
			"startDate": "2021-07-01T00:00:00+0000",
			"deadline": "2021-09-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"krIds": [
				"5fda249d289742000406b3e5",
				"5fdb72c63d2cf000035ceb38",
				"5ff445869f8c190003e0e445",
				"60bde9dc99b4177c88ad428c"
			],
			"childObjectiveIds": [],
			"latestUpdate": {
				"entityId": "5fdb72c63d2cf000035ceb37",
				"status": "UNDEFINED",
				"created": "2021-08-09T00:00:00+0000",
				"value": null,
				"description": ""
			},
			"periodAliasId": "602a6a2717378700039f342c",
			"weight": 0
		}
	],
	"krs": [
		{
			"id": "5fda249d289742000406b3e5",
			"key": "KR-2",
			"name": "adda",
			"description": null,
			"parentObjectiveId": "5fdb72c63d2cf000035ceb37",
			"issueIds": [
				"10003"
			],
			"ownerAccountId": "5c12ad9fd3af3b1ccfecbf55",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2020-12-16T15:15:41+0000",
			"startDate": "2021-01-05T00:00:00+0000",
			"deadline": "2021-03-28T00:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": null,
			"latestUpdate": {
				"entityId": "61c993aaa0fd9b768a0fb47d",
				"status": "AT_RISK",
				"created": "2021-12-27T10:21:30+0000",
				"value": 0.0,
				"description": null
			},
			"unit": null,
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 22.0,
				"jql": null
			},
			"weight": 1
		},
		{
			"id": "5fdb72c63d2cf000035ceb38",
			"key": "KR-3",
			"name": "12",
			"description": null,
			"parentObjectiveId": "5fdb72c63d2cf000035ceb37",
			"issueIds": [],
			"ownerAccountId": "5c12ad9fd3af3b1ccfecbf55",
			"collaboratorAccountIds": [
				"5c12ad9fd3af3b1ccfecbf55"
			],
			"percentDone": 0.0,
			"created": "2020-12-17T15:01:26+0000",
			"startDate": "2021-07-01T00:00:00+0000",
			"deadline": "2021-09-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342c",
			"latestUpdate": {
				"entityId": "5fdb72c63d2cf000035ceb39",
				"status": "ON_TRACK",
				"created": "2020-12-17T15:00:24+0000",
				"value": 1.0,
				"description": null
			},
			"unit": null,
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 1.0,
				"desiredValue": 2.0,
				"jql": null
			},
			"weight": 1
		},
		{
			"id": "5ff445869f8c190003e0e445",
			"key": "KR-5",
			"name": "miau",
			"description": null,
			"parentObjectiveId": "5fdb72c63d2cf000035ceb37",
			"issueIds": [],
			"ownerAccountId": "5c12ad9fd3af3b1ccfecbf55",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2021-01-05T10:55:02+0000",
			"startDate": "2021-07-01T00:00:00+0000",
			"deadline": "2021-09-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342c",
			"latestUpdate": {
				"entityId": "5ff445909f8c190003e0e447",
				"status": "DELAYED",
				"created": "2021-01-05T10:55:12+0000",
				"value": 0.0,
				"description": ""
			},
			"unit": null,
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 1
		},
		{
			"id": "605480b190c42b0003385170",
			"key": "KR-8",
			"name": "new auto KR",
			"description": null,
			"parentObjectiveId": "5fda249d289742000406b3e4",
			"issueIds": [
				"10005",
				"10006",
				"10010"
			],
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [
				"5c12ad9fd3af3b1ccfecbf55"
			],
			"percentDone": 33.33333333333333,
			"created": "2021-03-19T10:45:05+0000",
			"startDate": "2021-01-01T00:00:00+0000",
			"deadline": "2021-03-31T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342a",
			"latestUpdate": {
				"entityId": "620ea538512edb00acf67ac1",
				"status": "ON_TRACK",
				"created": "2022-02-17T19:42:43+0000",
				"value": 1.0,
				"description": null
			},
			"unit": null,
			"currentProgressDefinition": {
				"type": "AUTO",
				"startValue": null,
				"desiredValue": null,
				"jql": "project= TEST"
			},
			"weight": 1
		},
		{
			"id": "6061e921e2f4470003bc3210",
			"key": "KR-9",
			"name": "different start date",
			"description": null,
			"parentObjectiveId": "5fda249d289742000406b3e4",
			"issueIds": [
				"10000"
			],
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2021-03-29T14:50:09+0000",
			"startDate": "2020-04-01T00:00:00+0000",
			"deadline": "2020-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": null,
			"latestUpdate": {
				"entityId": "61138f2be5fd454858c3e1ee",
				"status": "AT_RISK",
				"created": "2021-08-11T08:49:47+0000",
				"value": 0.0,
				"description": null
			},
			"unit": {
				"name": "USD",
				"symbol": "$"
			},
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 1
		},
		{
			"id": "60bde9dc99b4177c88ad428c",
			"key": "KR-15",
			"name": "fghjfgh",
			"description": "",
			"parentObjectiveId": "5fdb72c63d2cf000035ceb37",
			"issueIds": [],
			"ownerAccountId": "5c12ad9fd3af3b1ccfecbf55",
			"collaboratorAccountIds": [
				"5c12ad9fd3af3b1ccfecbf55"
			],
			"percentDone": 0.0,
			"created": "2021-06-07T09:41:48+0000",
			"startDate": "2021-07-01T00:00:00+0000",
			"deadline": "2021-09-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [
				"605df77c6e53750003068c7d",
				"605dfb316e53750003068c8a"
			],
			"periodAliasId": "602a6a2717378700039f342c",
			"latestUpdate": {
				"entityId": "60bde9dc99b4177c88ad428d",
				"status": "ON_TRACK",
				"created": "2021-06-07T09:41:37+0000",
				"value": 0.0,
				"description": null
			},
			"unit": null,
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 1
		}
	],
	"teams": [
		{
			"id": "605df77c6e53750003068c7d",
			"name": "lets go!"
		},
		{
			"id": "605dfb316e53750003068c8a",
			"name": "My team"
		}
	],
	"periods": [
		{
			"id": "602a6a2717378700039f342a",
			"name": "Q1 Y2021",
			"startDate": "2021-01-01T00:00:00+0000",
			"deadline": "2021-03-31T23:59:59+0000"
		},
		{
			"id": "602a6a2717378700039f342c",
			"name": "Q3 Y2021",
			"startDate": "2021-07-01T00:00:00+0000",
			"deadline": "2021-09-30T23:59:59+0000"
		}
	],
	"labels": []
}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	ctx := context.Background()

	objectiveIDs := []string{"5fda249d289742000406b3e4", "5fdb72c63d2cf000035ceb37"}
	expand := []string{"OBJECTIVES", "KEY_RESULTS", "TEAMS", "PERIODS", "LABELS"}
	c := okrforjira.NewClient(client, token)
	got, err := c.ObjectivesByIDs(ctx, objectiveIDs, expand)
	assert.NoError(t, err)

	want := okrforjira.Response{
		OKRs: []okrforjira.OKR{
			{
				ID:                     "5fda249d289742000406b3e4",
				Key:                    "O-2",
				Name:                   "Become more mature company",
				Description:            "<p>This quarter we will be focusing on improving our performance.</p><p></p>",
				ParentObjectiveID:      "",
				OwnerAccountID:         "5c12ad9fd3af3b1ccfecbf55",
				CollaboratorAccountIDs: []string{},
				PercentDone:            8.333333333333332,
				Created:                time.Date(2020, time.December, 16, 15, 15, 41, 0, time.UTC),
				StartDate:              time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2021, time.March, 31, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				KRIDs: []string{
					"605480b190c42b0003385170",
					"6061e921e2f4470003bc3210",
				},
				ChildObjectiveIDs: []string{
					"5fdb72c63d2cf000035ceb37",
					"60743135b347480003dc6a9c",
					"61f9367df9aa7f0e4024a6fe",
				},
				LatestUpdate: okrforjira.Update{
					EntityID:    "5fda249d289742000406b3e4",
					Status:      "ON_TRACK",
					Created:     time.Date(2021, time.May, 5, 12, 15, 14, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				PeriodAliasID: "602a6a2717378700039f342a",
				Weight:        0.0,
			},
			{
				ID:                     "5fdb72c63d2cf000035ceb37",
				Key:                    "O-3",
				Name:                   "45",
				Description:            "",
				ParentObjectiveID:      "5fda249d289742000406b3e4",
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIDs: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2020, time.December, 17, 15, 1, 26, 0, time.UTC),
				StartDate:              time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2021, time.September, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				KRIDs: []string{
					"5fda249d289742000406b3e5",
					"5fdb72c63d2cf000035ceb38",
					"5ff445869f8c190003e0e445",
					"60bde9dc99b4177c88ad428c",
				},
				ChildObjectiveIDs: []string{},
				LatestUpdate: okrforjira.Update{
					EntityID:    "5fdb72c63d2cf000035ceb37",
					Status:      "UNDEFINED",
					Created:     time.Date(2021, time.August, 9, 0, 0, 0, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				PeriodAliasID: "602a6a2717378700039f342c",
				Weight:        0.0,
			},
		},
		KeyResults: []okrforjira.KeyResult{
			{
				ID:                "5fda249d289742000406b3e5",
				Key:               "KR-2",
				Name:              "adda",
				Description:       "",
				ParentObjectiveID: "5fdb72c63d2cf000035ceb37",
				IssueIDs: []string{
					"10003",
				},
				OwnerAccountID:         "5c12ad9fd3af3b1ccfecbf55",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2020, time.December, 16, 15, 15, 41, 0, time.UTC),
				StartDate:              time.Date(2021, time.January, 5, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2021, time.March, 28, 0, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "",
				LatestUpdate: okrforjira.Update{
					EntityID:    "61c993aaa0fd9b768a0fb47d",
					Status:      "AT_RISK",
					Created:     time.Date(2021, time.December, 27, 10, 21, 30, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 22.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
			{
				ID:                "5fdb72c63d2cf000035ceb38",
				Key:               "KR-3",
				Name:              "12",
				Description:       "",
				ParentObjectiveID: "5fdb72c63d2cf000035ceb37",
				IssueIDs:          []string{},
				OwnerAccountID:    "5c12ad9fd3af3b1ccfecbf55",
				CollaboratorAccountIds: []string{
					"5c12ad9fd3af3b1ccfecbf55",
				},
				PercentDone:   0.0,
				Created:       time.Date(2020, time.December, 17, 15, 1, 26, 0, time.UTC),
				StartDate:     time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				Deadline:      time.Date(2021, time.September, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:      []string{},
				TeamIDs:       []string{},
				PeriodAliasID: "602a6a2717378700039f342c",
				LatestUpdate: okrforjira.Update{
					EntityID:    "5fdb72c63d2cf000035ceb39",
					Status:      "ON_TRACK",
					Created:     time.Date(2020, time.December, 17, 15, 0, 24, 0, time.UTC),
					Value:       1.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   1.0,
					DesiredValue: 2.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
			{
				ID:                     "5ff445869f8c190003e0e445",
				Key:                    "KR-5",
				Name:                   "miau",
				Description:            "",
				ParentObjectiveID:      "5fdb72c63d2cf000035ceb37",
				IssueIDs:               []string{},
				OwnerAccountID:         "5c12ad9fd3af3b1ccfecbf55",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2021, time.January, 5, 10, 55, 02, 0, time.UTC),
				StartDate:              time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2021, time.September, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "602a6a2717378700039f342c",
				LatestUpdate: okrforjira.Update{
					EntityID:    "5ff445909f8c190003e0e447",
					Status:      "DELAYED",
					Created:     time.Date(2021, time.January, 5, 10, 55, 12, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
			{
				ID:                "605480b190c42b0003385170",
				Key:               "KR-8",
				Name:              "new auto KR",
				Description:       "",
				ParentObjectiveID: "5fda249d289742000406b3e4",
				IssueIDs: []string{
					"10005",
					"10006",
					"10010",
				},
				OwnerAccountID: "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIds: []string{
					"5c12ad9fd3af3b1ccfecbf55",
				},
				PercentDone:   33.33333333333333,
				Created:       time.Date(2021, time.March, 19, 10, 45, 5, 0, time.UTC),
				StartDate:     time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
				Deadline:      time.Date(2021, time.March, 31, 23, 59, 59, 0, time.UTC),
				LabelIDs:      []string{},
				TeamIDs:       []string{},
				PeriodAliasID: "602a6a2717378700039f342a",
				LatestUpdate: okrforjira.Update{
					EntityID:    "620ea538512edb00acf67ac1",
					Status:      "ON_TRACK",
					Created:     time.Date(2022, time.February, 17, 19, 42, 43, 0, time.UTC),
					Value:       1.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "AUTO",
					StartValue:   0.0,
					DesiredValue: 0.0,
					JQL:          "project= TEST",
				},
				Weight: 1.0,
			},
			{
				ID:                "6061e921e2f4470003bc3210",
				Key:               "KR-9",
				Name:              "different start date",
				Description:       "",
				ParentObjectiveID: "5fda249d289742000406b3e4",
				IssueIDs: []string{
					"10000",
				},
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2021, time.March, 29, 14, 50, 9, 0, time.UTC),
				StartDate:              time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2020, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "",
				LatestUpdate: okrforjira.Update{
					EntityID:    "61138f2be5fd454858c3e1ee",
					Status:      "AT_RISK",
					Created:     time.Date(2021, time.August, 11, 8, 49, 47, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "USD",
					Symbol: "$",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
			{
				ID:                "60bde9dc99b4177c88ad428c",
				Key:               "KR-15",
				Name:              "fghjfgh",
				Description:       "",
				ParentObjectiveID: "5fdb72c63d2cf000035ceb37",
				IssueIDs:          []string{},
				OwnerAccountID:    "5c12ad9fd3af3b1ccfecbf55",
				CollaboratorAccountIds: []string{
					"5c12ad9fd3af3b1ccfecbf55",
				},
				PercentDone: 0.0,
				Created:     time.Date(2021, time.June, 7, 9, 41, 48, 0, time.UTC),
				StartDate:   time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				Deadline:    time.Date(2021, time.September, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:    []string{},
				TeamIDs: []string{
					"605df77c6e53750003068c7d",
					"605dfb316e53750003068c8a",
				},
				PeriodAliasID: "602a6a2717378700039f342c",
				LatestUpdate: okrforjira.Update{
					EntityID:    "60bde9dc99b4177c88ad428d",
					Status:      "ON_TRACK",
					Created:     time.Date(2021, time.June, 7, 9, 41, 37, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
		},
		Teams: []okrforjira.Team{
			{
				ID:   "605df77c6e53750003068c7d",
				Name: "lets go!",
			},
			{
				ID:   "605dfb316e53750003068c8a",
				Name: "My team",
			},
		},
		Periods: []okrforjira.Period{
			{
				ID:        "602a6a2717378700039f342a",
				Name:      "Q1 Y2021",
				StartDate: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
				Deadline:  time.Date(2021, time.March, 31, 23, 59, 59, 0, time.UTC),
			},
			{
				ID:        "602a6a2717378700039f342c",
				Name:      "Q3 Y2021",
				StartDate: time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				Deadline:  time.Date(2021, time.September, 30, 23, 59, 59, 0, time.UTC),
			},
		},
		Labels: []okrforjira.Label{},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("unexpected result:\n%s", cmp.Diff(got, want))
	}
}

func TestClient_KeyResultsByDate(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "https://okr-for-jira-prod.herokuapp.com/api/v2/api-export/keyResults/byDate?startDateEpochMilli=1509412400000&deadlineEpochMilli=1748647410000&expand=OBJECTIVES,KEY_RESULTS,TEAMS,PERIODS,LABELS")
		assert.Equal(t, req.Header.Get("API-Token"), token)
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
	"okrs": [
		{
			"id": "620265fe81300a4c89e89806",
			"key": "O-99",
			"name": "Department objective ",
			"description": "",
			"parentObjectiveId": "620265d481300a4c89e89802",
			"ownerAccountId": "557058:a63fcd57-682a-450b-8d7c-ec330b2aa543",
			"collaboratorAccountIds": [],
			"percentDone": 4.761904761904762,
			"created": "2022-02-08T12:45:50+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [
				"61fac8f89556e8621145011d"
			],
			"krIds": [
				"62027c6f81300a4c89e89a1c"
			],
			"childObjectiveIds": [
				"6202664f81300a4c89e89810",
				"6202666281300a4c89e89812"
			],
			"latestUpdate": {
				"entityId": "620265fe81300a4c89e89806",
				"status": "NOT_STARTED",
				"created": "2022-02-08T12:45:50+0000",
				"value": null,
				"description": "Objective created."
			},
			"periodAliasId": "602a6a2717378700039f342f",
			"weight": 1
		},
		{
			"id": "6202664f81300a4c89e89810",
			"key": "O-101",
			"name": "Team objective",
			"description": "",
			"parentObjectiveId": "620265fe81300a4c89e89806",
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2022-02-08T12:47:11+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [
				"61c3402d8a875b727286a6c2"
			],
			"krIds": [
				"6202671581300a4c89e89827",
				"6202672781300a4c89e8982a"
			],
			"childObjectiveIds": [],
			"latestUpdate": {
				"entityId": "6202664f81300a4c89e89810",
				"status": "NOT_STARTED",
				"created": "2022-02-08T12:47:11+0000",
				"value": null,
				"description": "Objective created."
			},
			"periodAliasId": "602a6a2717378700039f342f",
			"weight": 1
		},
		{
			"id": "6202666281300a4c89e89812",
			"key": "O-102",
			"name": "Team 2 objective",
			"description": "",
			"parentObjectiveId": "620265fe81300a4c89e89806",
			"ownerAccountId": "5c12ad9fd3af3b1ccfecbf55",
			"collaboratorAccountIds": [],
			"percentDone": 9.523809523809524,
			"created": "2022-02-08T12:47:30+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [
				"611a1b88c385f85c13ac8630"
			],
			"krIds": [
				"6202673e81300a4c89e8982e",
				"6202674e81300a4c89e89830",
				"62029d8881300a4c89e8af27"
			],
			"childObjectiveIds": [],
			"latestUpdate": {
				"entityId": "6202666281300a4c89e89812",
				"status": "NOT_STARTED",
				"created": "2022-02-08T12:47:30+0000",
				"value": null,
				"description": "Objective created."
			},
			"periodAliasId": "602a6a2717378700039f342f",
			"weight": 1
		}
	],
	"krs": [
		{
			"id": "6202671581300a4c89e89827",
			"key": "KR-96",
			"name": "Team KR",
			"description": "",
			"parentObjectiveId": "6202664f81300a4c89e89810",
			"issueIds": [],
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2022-02-08T12:50:29+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342f",
			"latestUpdate": {
				"entityId": "6202671581300a4c89e89828",
				"status": "NOT_STARTED",
				"created": "2022-02-08T12:50:29+0000",
				"value": 0.0,
				"description": "Key result created."
			},
			"unit": {
				"name": "Numeric",
				"symbol": ""
			},
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 1
		},
		{
			"id": "6202672781300a4c89e8982a",
			"key": "KR-97",
			"name": "Team KR 2 ",
			"description": "",
			"parentObjectiveId": "6202664f81300a4c89e89810",
			"issueIds": [],
			"ownerAccountId": "557058:a63fcd57-682a-450b-8d7c-ec330b2aa543",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2022-02-08T12:50:47+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342f",
			"latestUpdate": {
				"entityId": "6202672781300a4c89e8982b",
				"status": "NOT_STARTED",
				"created": "2022-02-08T12:50:47+0000",
				"value": 0.0,
				"description": "Key result created."
			},
			"unit": {
				"name": "Numeric",
				"symbol": ""
			},
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 1
		},
		{
			"id": "6202673e81300a4c89e8982e",
			"key": "KR-98",
			"name": "Team KR 3",
			"description": "",
			"parentObjectiveId": "6202666281300a4c89e89812",
			"issueIds": [],
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2022-02-08T12:51:09+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342f",
			"latestUpdate": {
				"entityId": "6202673e81300a4c89e8982f",
				"status": "NOT_STARTED",
				"created": "2022-02-08T12:51:09+0000",
				"value": 0.0,
				"description": "Key result created."
			},
			"unit": {
				"name": "Numeric",
				"symbol": ""
			},
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 1
		},
		{
			"id": "6202674e81300a4c89e89830",
			"key": "KR-99",
			"name": "Team KR 4",
			"description": "",
			"parentObjectiveId": "6202666281300a4c89e89812",
			"issueIds": [],
			"ownerAccountId": "557058:a63fcd57-682a-450b-8d7c-ec330b2aa543",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2022-02-08T12:51:26+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342f",
			"latestUpdate": {
				"entityId": "6202674e81300a4c89e89831",
				"status": "NOT_STARTED",
				"created": "2022-02-08T12:51:26+0000",
				"value": 0.0,
				"description": "Key result created."
			},
			"unit": {
				"name": "Numeric",
				"symbol": ""
			},
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 1
		},
		{
			"id": "62027c6f81300a4c89e89a1c",
			"key": "KR-103",
			"name": "For tracking ",
			"description": "",
			"parentObjectiveId": "620265fe81300a4c89e89806",
			"issueIds": [],
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2022-02-08T14:21:35+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342f",
			"latestUpdate": {
				"entityId": "62027c6f81300a4c89e89a1d",
				"status": "NOT_STARTED",
				"created": "2022-02-08T14:21:35+0000",
				"value": 0.0,
				"description": "Key result created."
			},
			"unit": {
				"name": "Numeric",
				"symbol": ""
			},
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 1.0,
				"jql": null
			},
			"weight": 0
		},
		{
			"id": "62029d8881300a4c89e8af27",
			"key": "KR-104",
			"name": "Auto KR",
			"description": "",
			"parentObjectiveId": "6202666281300a4c89e89812",
			"issueIds": [
				"10008",
				"10007",
				"10009",
				"10000",
				"10002",
				"10001",
				"10004"
			],
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 28.57142857142857,
			"created": "2022-02-08T16:42:48+0000",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": "602a6a2717378700039f342f",
			"latestUpdate": {
				"entityId": "62029d8881300a4c89e8af2a",
				"status": "NOT_STARTED",
				"created": "2022-02-08T16:42:48+0000",
				"value": 2.0,
				"description": null
			},
			"unit": null,
			"currentProgressDefinition": {
				"type": "AUTO",
				"startValue": null,
				"desiredValue": null,
				"jql": "project = ddd"
			},
			"weight": 1
		}
	],
	"teams": [
		{
			"id": "611a1b88c385f85c13ac8630",
			"name": "Product team"
		},
		{
			"id": "61c3402d8a875b727286a6c2",
			"name": "Operations"
		},
		{
			"id": "61fac8f89556e8621145011d",
			"name": "Research"
		}
	],
	"periods": [
		{
			"id": "602a6a2717378700039f342f",
			"name": "Q2 Y2022",
			"startDate": "2022-04-01T00:00:00+0000",
			"deadline": "2022-06-30T23:59:59+0000"
		}
	],
	"labels": []
}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	ctx := context.Background()

	// Tuesday, October 31, 2017 3:13:20 PM GMT+01:00
	loc, _ := time.LoadLocation("Europe/Paris")
	startDate := time.Date(2017, time.October, 31, 2, 13, 20, 0, loc)
	deadline := time.Date(2025, time.May, 31, 1, 23, 30, 0, loc)
	expand := []string{"OBJECTIVES", "KEY_RESULTS", "TEAMS", "PERIODS", "LABELS"}
	c := okrforjira.NewClient(client, token)
	got, err := c.KeyResultsByDate(ctx, startDate, deadline, expand)
	assert.NoError(t, err)

	want := okrforjira.Response{
		OKRs: []okrforjira.OKR{
			{
				ID:                     "620265fe81300a4c89e89806",
				Key:                    "O-99",
				Name:                   "Department objective ",
				Description:            "",
				ParentObjectiveID:      "620265d481300a4c89e89802",
				OwnerAccountID:         "557058:a63fcd57-682a-450b-8d7c-ec330b2aa543",
				CollaboratorAccountIDs: []string{},
				PercentDone:            4.761904761904762,
				Created:                time.Date(2022, time.February, 8, 12, 45, 50, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs: []string{
					"61fac8f89556e8621145011d",
				},
				KRIDs: []string{
					"62027c6f81300a4c89e89a1c",
				},
				ChildObjectiveIDs: []string{
					"6202664f81300a4c89e89810",
					"6202666281300a4c89e89812",
				},
				LatestUpdate: okrforjira.Update{
					EntityID:    "620265fe81300a4c89e89806",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 12, 45, 50, 0, time.UTC),
					Value:       0.0,
					Description: "Objective created.",
				},
				PeriodAliasID: "602a6a2717378700039f342f",
				Weight:        1.0,
			},
			{
				ID:                     "6202664f81300a4c89e89810",
				Key:                    "O-101",
				Name:                   "Team objective",
				Description:            "",
				ParentObjectiveID:      "620265fe81300a4c89e89806",
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIDs: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2022, time.February, 8, 12, 47, 11, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs: []string{
					"61c3402d8a875b727286a6c2",
				},
				KRIDs: []string{
					"6202671581300a4c89e89827",
					"6202672781300a4c89e8982a",
				},
				ChildObjectiveIDs: []string{},
				LatestUpdate: okrforjira.Update{
					EntityID:    "6202664f81300a4c89e89810",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 12, 47, 11, 0, time.UTC),
					Value:       0.0,
					Description: "Objective created.",
				},
				PeriodAliasID: "602a6a2717378700039f342f",
				Weight:        1.0,
			},
			{
				ID:                     "6202666281300a4c89e89812",
				Key:                    "O-102",
				Name:                   "Team 2 objective",
				Description:            "",
				ParentObjectiveID:      "620265fe81300a4c89e89806",
				OwnerAccountID:         "5c12ad9fd3af3b1ccfecbf55",
				CollaboratorAccountIDs: []string{},
				PercentDone:            9.523809523809524,
				Created:                time.Date(2022, time.February, 8, 12, 47, 30, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs: []string{
					"611a1b88c385f85c13ac8630",
				},
				KRIDs: []string{
					"6202673e81300a4c89e8982e",
					"6202674e81300a4c89e89830",
					"62029d8881300a4c89e8af27",
				},
				ChildObjectiveIDs: []string{},
				LatestUpdate: okrforjira.Update{
					EntityID:    "6202666281300a4c89e89812",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 12, 47, 30, 0, time.UTC),
					Value:       0.0,
					Description: "Objective created.",
				},
				PeriodAliasID: "602a6a2717378700039f342f",
				Weight:        1.0,
			},
		},
		KeyResults: []okrforjira.KeyResult{
			{
				ID:                     "6202671581300a4c89e89827",
				Key:                    "KR-96",
				Name:                   "Team KR",
				Description:            "",
				ParentObjectiveID:      "6202664f81300a4c89e89810",
				IssueIDs:               []string{},
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2022, time.February, 8, 12, 50, 29, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "602a6a2717378700039f342f",
				LatestUpdate: okrforjira.Update{
					EntityID:    "6202671581300a4c89e89828",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 12, 50, 29, 0, time.UTC),
					Value:       0.0,
					Description: "Key result created.",
				},
				Unit: okrforjira.Unit{
					Name:   "Numeric",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
			{
				ID:                     "6202672781300a4c89e8982a",
				Key:                    "KR-97",
				Name:                   "Team KR 2 ",
				Description:            "",
				ParentObjectiveID:      "6202664f81300a4c89e89810",
				IssueIDs:               []string{},
				OwnerAccountID:         "557058:a63fcd57-682a-450b-8d7c-ec330b2aa543",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2022, time.February, 8, 12, 50, 47, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "602a6a2717378700039f342f",
				LatestUpdate: okrforjira.Update{
					EntityID:    "6202672781300a4c89e8982b",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 12, 50, 47, 0, time.UTC),
					Value:       0.0,
					Description: "Key result created.",
				},
				Unit: okrforjira.Unit{
					Name:   "Numeric",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
			{
				ID:                     "6202673e81300a4c89e8982e",
				Key:                    "KR-98",
				Name:                   "Team KR 3",
				Description:            "",
				ParentObjectiveID:      "6202666281300a4c89e89812",
				IssueIDs:               []string{},
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2022, time.February, 8, 12, 51, 9, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "602a6a2717378700039f342f",
				LatestUpdate: okrforjira.Update{
					EntityID:    "6202673e81300a4c89e8982f",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 12, 51, 9, 0, time.UTC),
					Value:       0.0,
					Description: "Key result created.",
				},
				Unit: okrforjira.Unit{
					Name:   "Numeric",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
			{
				ID:                     "6202674e81300a4c89e89830",
				Key:                    "KR-99",
				Name:                   "Team KR 4",
				Description:            "",
				ParentObjectiveID:      "6202666281300a4c89e89812",
				IssueIDs:               []string{},
				OwnerAccountID:         "557058:a63fcd57-682a-450b-8d7c-ec330b2aa543",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2022, time.February, 8, 12, 51, 26, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "602a6a2717378700039f342f",
				LatestUpdate: okrforjira.Update{
					EntityID:    "6202674e81300a4c89e89831",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 12, 51, 26, 0, time.UTC),
					Value:       0.0,
					Description: "Key result created.",
				},
				Unit: okrforjira.Unit{
					Name:   "Numeric",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
			{
				ID:                     "62027c6f81300a4c89e89a1c",
				Key:                    "KR-103",
				Name:                   "For tracking ",
				Description:            "",
				ParentObjectiveID:      "620265fe81300a4c89e89806",
				IssueIDs:               []string{},
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2022, time.February, 8, 14, 21, 35, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "602a6a2717378700039f342f",
				LatestUpdate: okrforjira.Update{
					EntityID:    "62027c6f81300a4c89e89a1d",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 14, 21, 35, 0, time.UTC),
					Value:       0.0,
					Description: "Key result created.",
				},
				Unit: okrforjira.Unit{
					Name:   "Numeric",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 1.0,
					JQL:          "",
				},
				Weight: 0.0,
			},
			{
				ID:                "62029d8881300a4c89e8af27",
				Key:               "KR-104",
				Name:              "Auto KR",
				Description:       "",
				ParentObjectiveID: "6202666281300a4c89e89812",
				IssueIDs: []string{
					"10008",
					"10007",
					"10009",
					"10000",
					"10002",
					"10001",
					"10004",
				},
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIds: []string{},
				PercentDone:            28.57142857142857,
				Created:                time.Date(2022, time.February, 8, 16, 42, 48, 0, time.UTC),
				StartDate:              time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "602a6a2717378700039f342f",
				LatestUpdate: okrforjira.Update{
					EntityID:    "62029d8881300a4c89e8af2a",
					Status:      "NOT_STARTED",
					Created:     time.Date(2022, time.February, 8, 16, 42, 48, 0, time.UTC),
					Value:       2.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "AUTO",
					StartValue:   0.0,
					DesiredValue: 0.0,
					JQL:          "project = ddd",
				},
				Weight: 1.0,
			},
		},
		Teams: []okrforjira.Team{
			{
				ID:   "611a1b88c385f85c13ac8630",
				Name: "Product team",
			},
			{
				ID:   "61c3402d8a875b727286a6c2",
				Name: "Operations",
			},
			{
				ID:   "61fac8f89556e8621145011d",
				Name: "Research",
			},
		},
		Periods: []okrforjira.Period{
			{
				ID:        "602a6a2717378700039f342f",
				Name:      "Q2 Y2022",
				StartDate: time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
				Deadline:  time.Date(2022, time.June, 30, 23, 59, 59, 0, time.UTC),
			},
		},
		Labels: []okrforjira.Label{},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("unexpected result:\n%s", cmp.Diff(got, want))
	}
}

func TestClient_KeyResultsByIDs(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "https://okr-for-jira-prod.herokuapp.com/api/v2/api-export/keyResults/byIds?keyResultIds=5fda249d289742000406b3e5&expand=OBJECTIVES,KEY_RESULTS,TEAMS,PERIODS,LABELS")
		assert.Equal(t, req.Header.Get("API-Token"), token)
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
	"okrs": [
		{
			"id": "5fdb72c63d2cf000035ceb37",
			"key": "O-3",
			"name": "45",
			"description": null,
			"parentObjectiveId": "5fda249d289742000406b3e4",
			"ownerAccountId": "5dbfee8570f1ea0df7698353",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2020-12-17T15:01:26+0000",
			"startDate": "2021-07-01T00:00:00+0000",
			"deadline": "2021-09-30T23:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"krIds": [
				"5fda249d289742000406b3e5",
				"5fdb72c63d2cf000035ceb38",
				"5ff445869f8c190003e0e445",
				"60bde9dc99b4177c88ad428c"
			],
			"childObjectiveIds": [],
			"latestUpdate": {
				"entityId": "5fdb72c63d2cf000035ceb37",
				"status": "UNDEFINED",
				"created": "2021-08-09T00:00:00+0000",
				"value": null,
				"description": ""
			},
			"periodAliasId": "602a6a2717378700039f342c",
			"weight": 0
		}
	],
	"krs": [
		{
			"id": "5fda249d289742000406b3e5",
			"key": "KR-2",
			"name": "adda",
			"description": null,
			"parentObjectiveId": "5fdb72c63d2cf000035ceb37",
			"issueIds": [
				"10003"
			],
			"ownerAccountId": "5c12ad9fd3af3b1ccfecbf55",
			"collaboratorAccountIds": [],
			"percentDone": 0.0,
			"created": "2020-12-16T15:15:41+0000",
			"startDate": "2021-01-05T00:00:00+0000",
			"deadline": "2021-03-28T00:59:59+0000",
			"labelIds": [],
			"teamIds": [],
			"periodAliasId": null,
			"latestUpdate": {
				"entityId": "61c993aaa0fd9b768a0fb47d",
				"status": "AT_RISK",
				"created": "2021-12-27T10:21:30+0000",
				"value": 0.0,
				"description": null
			},
			"unit": null,
			"currentProgressDefinition": {
				"type": "STANDARD",
				"startValue": 0.0,
				"desiredValue": 22.0,
				"jql": null
			},
			"weight": 1
		}
	],
	"teams": [],
	"periods": [
		{
			"id": "602a6a2717378700039f342c",
			"name": "Q3 Y2021",
			"startDate": "2021-07-01T00:00:00+0000",
			"deadline": "2021-09-30T23:59:59+0000"
		}
	],
	"labels": []
}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	ctx := context.Background()

	keyResultIDs := []string{"5fda249d289742000406b3e5"}
	expand := []string{"OBJECTIVES", "KEY_RESULTS", "TEAMS", "PERIODS", "LABELS"}
	c := okrforjira.NewClient(client, token)
	got, err := c.KeyResultsByIDs(ctx, keyResultIDs, expand)
	assert.NoError(t, err)

	want := okrforjira.Response{
		OKRs: []okrforjira.OKR{
			{
				ID:                     "5fdb72c63d2cf000035ceb37",
				Key:                    "O-3",
				Name:                   "45",
				Description:            "",
				ParentObjectiveID:      "5fda249d289742000406b3e4",
				OwnerAccountID:         "5dbfee8570f1ea0df7698353",
				CollaboratorAccountIDs: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2020, time.December, 17, 15, 1, 26, 0, time.UTC),
				StartDate:              time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2021, time.September, 30, 23, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				KRIDs: []string{
					"5fda249d289742000406b3e5",
					"5fdb72c63d2cf000035ceb38",
					"5ff445869f8c190003e0e445",
					"60bde9dc99b4177c88ad428c",
				},
				ChildObjectiveIDs: []string{},
				LatestUpdate: okrforjira.Update{
					EntityID:    "5fdb72c63d2cf000035ceb37",
					Status:      "UNDEFINED",
					Created:     time.Date(2021, time.August, 9, 0, 0, 0, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				PeriodAliasID: "602a6a2717378700039f342c",
				Weight:        0.0,
			},
		},
		KeyResults: []okrforjira.KeyResult{
			{
				ID:                     "5fda249d289742000406b3e5",
				Key:                    "KR-2",
				Name:                   "adda",
				Description:            "",
				ParentObjectiveID:      "5fdb72c63d2cf000035ceb37",
				IssueIDs:               []string{"10003"},
				OwnerAccountID:         "5c12ad9fd3af3b1ccfecbf55",
				CollaboratorAccountIds: []string{},
				PercentDone:            0.0,
				Created:                time.Date(2020, time.December, 16, 15, 15, 41, 0, time.UTC),
				StartDate:              time.Date(2021, time.January, 5, 0, 0, 0, 0, time.UTC),
				Deadline:               time.Date(2021, time.March, 28, 0, 59, 59, 0, time.UTC),
				LabelIDs:               []string{},
				TeamIDs:                []string{},
				PeriodAliasID:          "",
				LatestUpdate: okrforjira.Update{
					EntityID:    "61c993aaa0fd9b768a0fb47d",
					Status:      "AT_RISK",
					Created:     time.Date(2021, time.December, 27, 10, 21, 30, 0, time.UTC),
					Value:       0.0,
					Description: "",
				},
				Unit: okrforjira.Unit{
					Name:   "",
					Symbol: "",
				},
				CurrentProgressDefinition: okrforjira.ProgressDefinition{
					Type:         "STANDARD",
					StartValue:   0.0,
					DesiredValue: 22.0,
					JQL:          "",
				},
				Weight: 1.0,
			},
		},
		Teams: []okrforjira.Team{},
		Periods: []okrforjira.Period{
			{
				ID:        "602a6a2717378700039f342c",
				Name:      "Q3 Y2021",
				StartDate: time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				Deadline:  time.Date(2021, time.September, 30, 23, 59, 59, 0, time.UTC),
			},
		},
		Labels: []okrforjira.Label{},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("unexpected result:\n%s", cmp.Diff(got, want))
	}
}

func TestClient_UpdateObjective(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "https://okr-for-jira-prod.herokuapp.com/api/v2/api-update/objectives")
		assert.Equal(t, req.Header.Get("API-Token"), token)
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
	"entityId": "62334eac00ee2b102e34fdb7",
	"status": "ON TRACK",
	"created": "2022-05-20T09:58:09+0000",
	"value": null,
	"description": "Spaceship assembly docks are delivering on time"
}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	ctx := context.Background()

	const (
		objectiveID = "62334eac00ee2b102e34fdb7"
		status      = "ON TRACK"
		description = "Spaceship assembly docks are delivering on time"
	)
	c := okrforjira.NewClient(client, token)
	got, err := c.UpdateObjective(ctx, objectiveID, status, description)
	assert.NoError(t, err)

	want := okrforjira.Update{
		EntityID:    objectiveID,
		Status:      status,
		Created:     time.Date(2022, time.May, 20, 9, 58, 9, 0, time.UTC),
		Value:       0.0,
		Description: description,
	}
	if !cmp.Equal(got, want) {
		t.Errorf("unexpected result:\n%s", cmp.Diff(got, want))
	}
}

func TestClient_UpdateKeyResult(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "https://okr-for-jira-prod.herokuapp.com/api/v2/api-update/keyResults")
		assert.Equal(t, req.Header.Get("API-Token"), token)
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
	"entityId": "62384a6942adda046598b3bd",
	"status": "AT RISK",
	"created": "2022-05-20T13:01:35+0000",
	"value": 13500.5,
	"description": "Reduction in ship hull output is caused by Unobtainium supply disruptions."
}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	ctx := context.Background()

	const (
		keyResultID = "62384a6942adda046598b3bd"
		status      = "AT RISK"
		value       = 13500.5
		description = "Reduction in ship hull output is caused by Unobtainium supply disruptions."
	)
	c := okrforjira.NewClient(client, token)
	got, err := c.UpdateKeyResult(ctx, keyResultID, status, value, description)
	assert.NoError(t, err)

	want := okrforjira.Update{
		EntityID:    keyResultID,
		Status:      status,
		Created:     time.Date(2022, time.May, 20, 13, 1, 35, 0, time.UTC),
		Value:       value,
		Description: description,
	}
	if !cmp.Equal(got, want) {
		t.Errorf("unexpected result:\n%s", cmp.Diff(got, want))
	}
}

// Instead of the default http.Transport we will use http.RoundTripper.
// It will allow us to fake the server.

// RoundTripFunc is a function type that implements the RoundTripper interface.
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
