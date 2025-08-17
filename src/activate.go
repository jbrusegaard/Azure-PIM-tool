package src

import (
	"app/azuClient"
	"app/constants"
	"fmt"
)

type ActivationOptions struct {
	Reason    string
	Duration  int    // Duration in hours
	GroupName string // Filter criteria for activation
}

func ActivatePim(opts ActivationOptions) {
	appSettings := Initialize()
	azureClient := azuClient.AzureClient{
		AzurePimToken: appSettings.Session.AZPimToken,
	}
	eligibleRoleMap, err := azureClient.GetEligibleRoles(constants.AzurePimGroupApiUrlRoleAssignments)

	if len(eligibleRoleMap) == 0 {
		fmt.Println("No eligible roles found.")
	}
	roleToActivate, found := eligibleRoleMap[opts.GroupName]

	if !found {
		fmt.Println("No eligible group found with the specified name:", opts.GroupName)
		return
	}

	fmt.Println("Activating role:", roleToActivate.GetGroupName())
	requestBody := azuClient.BuildPimRequestBody(
		roleToActivate,
		roleToActivate.Subject.Id,
		opts.Reason,
		opts.Duration,
	)
	resp, err := azureClient.Activate(constants.AzurePimGroupApiUrlRoleAssigmentRequests, requestBody)
	if err != nil {
		panic("Failed to activate role: " + err.Error())
	}
	fmt.Println("Activation response:", resp)

	fmt.Println("YAY WE DID IT!!!")

}
