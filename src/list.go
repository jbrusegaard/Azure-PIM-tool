package src

import (
	"app/azuClient"
	"app/constants"
	"app/log"
)

type ListOpts struct {
	Interactive bool
	Headless    bool
	Deubg       bool
}

func ListGroups(opts ListOpts) {
	logger := log.InitializeLogger()

	// Set debugging
	setDebugging(opts.Deubg)

	appSettings := Initialize(logger, InitOpts{
		Interactive: opts.Interactive,
		Headless:    opts.Headless,
	})
	client := azuClient.AzureClient{
		AzurePimToken: appSettings.Session.AZPimToken,
	}
	roles, err := client.GetEligibleRoles(constants.AzurePimGroupApiUrlRoleAssignments)
	if err != nil {
		exitWithError(logger, "Error fetching eligible roles", err.Error())
	}
	if len(roles) == 0 {
		logger.Warn("No eligible roles found")
		return
	}

	for _, role := range roles {
		logger.WithPrefix("ROLE").Info(role.GetGroupName())
	}
}
