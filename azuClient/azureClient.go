package azuClient

type AzureClient struct {
	AzurePimToken AzurePimToken
}

func (a *AzureClient) GetToken() AzurePimToken {
	return a.AzurePimToken
}

func (a *AzureClient) SetToken(token AzurePimToken) {
	a.AzurePimToken = token
}

func (a *AzureClient) MakePIMRequest(url string, role string) (string, error) {
	return "", nil // Placeholder for actual implementation
}
