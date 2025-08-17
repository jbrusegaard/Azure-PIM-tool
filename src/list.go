package src

import (
	"app/azuClient"
	"app/constants"
	"fmt"
)

func ListGroups() {
	appSettings := Initialize()
	client := azuClient.AzureClient{
		AzurePimToken: appSettings.Session.AZPimToken,
	}
	roles, err := client.GetEligibleRoles(constants.AzurePimGroupApiUrlRoleAssignments)
	if err != nil {
		panic(err)
	}
	for _, role := range roles {
		fmt.Println(role.GetGroupName())
	}
}
