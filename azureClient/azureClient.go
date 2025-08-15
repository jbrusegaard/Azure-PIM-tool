package azureclient

type AzureClient struct {
	AccessToken string
	ObjectID    string
	Email       string
}

func (a *AzureClient) GetAccessToken() string {
	return a.AccessToken
}
