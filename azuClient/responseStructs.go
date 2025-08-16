package azuClient

import "fmt"

type AzureGroupResponseSubject struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	DisplayName   string `json:"displayName"`
	Email         string `json:"email"`
	PrincipalName string `json:"principalName"`
}
type AzureGroupResponseRoleDefinitionResource struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	ExternalId  string `json:"externalId"`
	Status      string `json:"status"`
}
type AzureGroupResponseRoleDefinition struct {
	Id         string                                   `json:"id"`
	ResourceId string                                   `json:"resourceId"`
	ExternalId string                                   `json:"externalId"`
	Resource   AzureGroupResponseRoleDefinitionResource `json:"resource"`
}

type AzureGroupResponse struct {
	Id               string                           `json:"id"`
	ResourceId       string                           `json:"resourceId"`
	RoleDefinitionId string                           `json:"roleDefinitionId"`
	Subject          AzureGroupResponseSubject        `json:"subject"`
	RoleDefinition   AzureGroupResponseRoleDefinition `json:"roleDefinition"`
}

type AzureGroupResponseList struct {
	Value []AzureGroupResponse `json:"value"`
}

func DisplayEligibleRoles(roles AzureGroupResponseList) string {
	if len(roles.Value) == 0 {
		return "No eligible roles found."
	}

	result := "Eligible Roles:\n"
	for _, role := range roles.Value {
		result += fmt.Sprintf("%s\n", role.RoleDefinition.Resource.DisplayName)
	}
	return result
}

func ComputeEligibleRoles(roles AzureGroupResponseList) []AzureGroupResponse {
	var eligibleRoles []AzureGroupResponse
	for _, role := range roles.Value {
		if role.Subject.Id != "" && role.RoleDefinition.Resource.DisplayName != "" {
			eligibleRoles = append(eligibleRoles, role)
		}
	}
	return eligibleRoles
}
