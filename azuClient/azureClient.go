package azuClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AzureClient struct {
	AzurePimToken AzurePimToken
}

type PimRequestSchedule struct {
	Type          string  `json:"type"`
	StartDateTime *string `json:"startDateTime"`
	EndDateTime   *string `json:"endDateTime"`
	Duration      string  `json:"duration"`
}

type AzurePimRequestBody struct {
	RoleDefinitionID               string             `json:"roleDefinitionId"`
	ResourceID                     string             `json:"resourceId"`
	SubjectID                      string             `json:"subjectId"`
	AssignmentState                string             `json:"assignmentState"`
	Type                           string             `json:"type"`
	Reason                         string             `json:"reason"`
	TicketNumber                   string             `json:"ticketNumber"`
	TicketSystem                   string             `json:"ticketSystem"`
	ScopedResourceID               string             `json:"scopedResourceId"`
	LinkedEligibleRoleAssignmentID string             `json:"linkedEligibleRoleAssignmentId"`
	Schedule                       PimRequestSchedule `json:"schedule"`
}

type ODataSearchParams struct {
	Filter string `url:"$filter"`
	Expand string `url:"$expand"`
}

func getPimRequestSchedule(durationHours int) PimRequestSchedule {
	durationMinutes := durationHours * 60
	return PimRequestSchedule{
		Type:          "Once",
		StartDateTime: nil,
		EndDateTime:   nil,
		Duration:      fmt.Sprintf("PT%dM", durationMinutes),
	}
}

func BuildPimRequestBody(
	role AzureGroupResponse, subjectID string, reason string, durationHours int,
) AzurePimRequestBody {
	return AzurePimRequestBody{
		RoleDefinitionID:               role.RoleDefinitionId,
		ResourceID:                     role.ResourceId,
		SubjectID:                      subjectID,
		Reason:                         reason,
		LinkedEligibleRoleAssignmentID: role.Id,
		Schedule:                       getPimRequestSchedule(durationHours),
		AssignmentState:                "Active",
		Type:                           "UserAdd",
		TicketNumber:                   "",
		TicketSystem:                   "",
		ScopedResourceID:               "",
	}
}

func (a *AzureClient) GetToken() AzurePimToken {
	return a.AzurePimToken
}

func (a *AzureClient) SetToken(token AzurePimToken) {
	a.AzurePimToken = token
}

func (a *AzureClient) GetEligibleRoles(base_url string) (string, error) {
	params := url.Values{}
	params.Add("$filter", "(subject/id eq '"+a.AzurePimToken.SubjectID+"') and (assignmentState eq 'Eligible')")
	params.Add("$expand", "linkedEligibleRoleAssignment,subject,scopedResource,roleDefinition($expand=resource)")
	req_url, err := url.Parse(base_url)
	req_url.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, req_url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + a.AzurePimToken.Secret,
		"Content-Type":  "application/json",
		"Accept":        "*/*",
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(resBody))
	}
	return string(resBody), nil
}

func (a *AzureClient) Activate(url string, body AzurePimRequestBody) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost, url, bytes.NewBuffer(payload),
	)
	headers := map[string]string{
		"Authorization": "Bearer " + a.AzurePimToken.Secret,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(resBody))
	}
	return string(resBody), err
}
