package azuClient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type AzurePimToken struct {
	CredentialType string `json:"credentialType"`
	Secret         string `json:"secret"`
	TokenType      string `json:"tokenType"`
	ClientID       string `json:"clientId"`
	Realm          string `json:"realm"`
	ExpiresOn      string `json:"expiresOn"`
	HomeAccountId  string `json:"homeAccountId"`
	Target         string `json:"target"`
	Environment    string `json:"environment"`
	SubjectID      string `json:"subjectId"`
	Email          string `json:"email"`
}

// AzureClaims represents the claims structure for Azure JWT tokens
type AzureClaims struct {
	ObjectID string `json:"oid"`   // Azure Object ID (may not always be present)
	Email    string `json:"email"` // User email
}

func (apt *AzurePimToken) ComputeAdditionalFields() error {
	if apt.SubjectID != "" {
		return nil
	}

	// Manually decode the JWT token to extract claims without signature verification
	// Split the JWT into parts (header.payload.signature)
	parts := strings.Split(apt.Secret, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	// Decode the payload (claims) part
	payload := parts[1]

	// Add padding if needed for base64 decoding
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	// Decode base64url to base64
	payload = strings.ReplaceAll(payload, "-", "+")
	payload = strings.ReplaceAll(payload, "_", "/")

	// Decode the base64 payload
	decodedPayload, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return fmt.Errorf("failed to decode JWT payload: %v", err)
	}
	// Parse into our custom claims structure
	var claims AzureClaims
	if err := json.Unmarshal(decodedPayload, &claims); err != nil {
		return fmt.Errorf("failed to unmarshal JWT claims: %v", err)
	}

	apt.SubjectID = claims.ObjectID
	apt.Email = claims.Email
	return nil
}
