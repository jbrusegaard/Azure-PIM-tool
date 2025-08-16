package azuClient

type AzurePimToken struct {
	CredentialType string `json:"credentialType"`
	Secret         string `json:"secret"`
	TokenType      string `json:"tokenType"`
	ClientID       string `json:"clientId"`
	Realm          string `json:"realm"`
	ExpiresOn      string `json:"expiresOn"`
}
