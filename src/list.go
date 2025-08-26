package src

import (
	"app/azuClient"
	"app/constants"
	"app/log"
)

type ListOpts struct {
	Interactive bool
	Headless    bool
}

func ListGroups(opts ListOpts) {
	logger := log.InitializeLogger()
	appSettings := Initialize(logger, InitOpts(opts))
	client := azuClient.AzureClient{
		AzurePimToken: appSettings.Session.AZPimToken,
	}
	roles, err := client.GetEligibleRoles(constants.AzurePimGroupApiUrlRoleAssignments)
	if err != nil {
		panic(err)
	}
	if len(roles) == 0 {
		logger.Warn("No eligible roles found")
		return
	}

	for _, role := range roles {
		logger.WithPrefix("ROLE").Info(role.GetGroupName())
	}
}
