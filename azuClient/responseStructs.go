package azuClient

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
