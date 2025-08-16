package src

import (
	"app/azuClient"
	"app/constants"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/playwright-community/playwright-go"
)

type ActivationOptions struct {
	Reason         string
	Duration       int    // Duration in hours
	ActivationType string // Type of activation, e.g., "group", "resource", "role"
	Filter         string // Filter criteria for activation
}

func preflight() {
	// check if playwright is installed
	if _, err := playwright.Run(); err != nil {
		iErr := playwright.Install(&playwright.RunOptions{Browsers: []string{"chromium"}})
		if iErr != nil {
			panic("Failed to install Playwright: " + iErr.Error())
		}
	}
}

func ActivatePim(opts ActivationOptions) {
	preflight()
	appSettings := Initialize()
	now := time.Now().Unix()
	expiresOn, err := strconv.Atoi(appSettings.Session.AZPimToken.ExpiresOn)
	if err != nil {
		expiresOn = 0
	}
	if now > int64(expiresOn) {
		fmt.Println("ActivatePim: Token expired")
		LaunchBrowserToGetToken(
			appSettings, PimOptions{
				Headless:        false,
				AppMode:         true,
				KioskMode:       true,
				PreserveSession: true,
				AzurePortalURL:  constants.AZURE_PORTAL_URL,
			},
		)
		appSettings = Initialize()
	}

	azureClient := azuClient.AzureClient{
		AzurePimToken: appSettings.Session.AZPimToken,
	}

	res, err := azureClient.GetEligibleRoles(constants.AZURE_PIM_GROUP_API_URL_ROLE_ASSIGNMENTS)
	if err != nil {
		panic("Failed to get eligible roles: " + err.Error())
	}
	var eligibleRoles azuClient.AzureGroupResponseList
	err = json.Unmarshal([]byte(res), &eligibleRoles)
	if err != nil {
		panic("Failed to unmarshal eligible roles: " + err.Error())
	}
	eligibleRoleMap := azuClient.ComputeEligibleRoles(eligibleRoles)
	if len(eligibleRoleMap) == 0 {
		fmt.Println("No eligible roles found.")
	}
	roleToActivate, found := eligibleRoleMap[opts.Filter]

	if !found {
		fmt.Println("No eligible group found with the specified filter:", opts.Filter)
		return
	}

	fmt.Println("Activating role:", roleToActivate.RoleDefinition.Resource.DisplayName)
	requestBody := azuClient.BuildPimRequestBody(
		roleToActivate,
		roleToActivate.Subject.Id,
		opts.Reason,
		opts.Duration,
	)
	resp, err := azureClient.Activate(constants.AZURE_PIM_GROUP_API_URL_ROLE_ASSIGMENT_REQUESTS, requestBody)
	if err != nil {
		panic("Failed to activate role: " + err.Error())
	}
	fmt.Println("Activation response:", resp)

	fmt.Println("YAY WE DID IT!!!")

}
