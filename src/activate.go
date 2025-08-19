package src

import (
	"encoding/json"

	"app/azuClient"
	"app/constants"
	"app/log"
)

type ActivationOptions struct {
	Headless    bool
	Interactive bool
	Reason      string
	Duration    int      // Duration in hours
	GroupNames  []string // Filter criteria for activation
}

func ActivatePim(opts ActivationOptions) {
	logger := log.InitializeLogger()
	appSettings := Initialize(logger, InitOpts{
		Interactive: opts.Interactive,
		Headless:    opts.Headless,
	})
	azureClient := azuClient.AzureClient{
		AzurePimToken: appSettings.Session.AZPimToken,
	}
	eligibleRoleMap, err := azureClient.GetEligibleRoles(constants.AzurePimGroupApiUrlRoleAssignments)
	if err != nil {
		logger.Errorf("Error fetching eligible roles: %s", err)
		return
	}

	if len(eligibleRoleMap) == 0 {
		logger.Warn("No eligible roles found.")
	}
	for _, groupName := range opts.GroupNames {
		roleToActivate, found := eligibleRoleMap[groupName]

		if !found {
			logger.Warnf("No eligible group found with the specified name: %s", groupName)
			continue
		}

		logger.With("role", roleToActivate.GetGroupName()).Info("Activating role")
		requestBody := azuClient.BuildPimRequestBody(
			roleToActivate,
			opts.Reason,
			opts.Duration,
		)
		resp, err := azureClient.Activate(constants.AzurePimGroupApiUrlRoleAssigmentRequests, requestBody)
		if err != nil {
			if resp != "" {
				var errResp *azuClient.AzureGroupErrorResponse
				unmarshErr := json.Unmarshal([]byte(resp), &errResp)
				if unmarshErr != nil {
					logger.With("role", groupName).Error(err.Error())
					return
				}
				logger.With("role", groupName).Warn(errResp.Error.Message)
			} else {
				logger.With("role", groupName).Error(err.Error())
			}
			return
		}
		logger.With("role", roleToActivate.GetGroupName()).Info("Successfully activated role")
	}

}
