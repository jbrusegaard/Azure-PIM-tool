package src

import (
	"encoding/json"

	"app/azuClient"
	"app/constants"

	"github.com/charmbracelet/log"
)

type ActivationOptions struct {
	Reason     string
	Duration   int      // Duration in hours
	GroupNames []string // Filter criteria for activation
}

func ActivatePim(opts ActivationOptions) {
	appSettings := Initialize()
	azureClient := azuClient.AzureClient{
		AzurePimToken: appSettings.Session.AZPimToken,
	}
	eligibleRoleMap, err := azureClient.GetEligibleRoles(constants.AzurePimGroupApiUrlRoleAssignments)
	if err != nil {
		log.Error("Error fetching eligible roles:", err)
		return
	}

	if len(eligibleRoleMap) == 0 {
		log.Warn("No eligible roles found.")
	}
	for _, groupName := range opts.GroupNames {
		roleToActivate, found := eligibleRoleMap[groupName]

		if !found {
			log.Warnf("No eligible group found with the specified name: %s", groupName)
			continue
		}

		log.Infof("Activating role: %s", roleToActivate.GetGroupName())
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
					log.Error(err.Error())
					return
				}
				log.Warn(errResp.Error.Message)
			} else {
				log.Error(err.Error())
			}
			return
		}
		log.Infof("Successfully activated role: %s!", roleToActivate.GetGroupName())
	}

}
