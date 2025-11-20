package azuClient

import (
	"app/log"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type AzurePimToken struct {
	Audience       string   `json:"audience"`
	ClientID       string   `json:"clientId"`
	CredentialType string   `json:"credentialType"`
	Email          string   `json:"email"`
	Environment    string   `json:"environment"`
	ExpiresOn      string   `json:"expiresOn"`
	Groups         []string `json:"groups"`
	HomeAccountId  string   `json:"homeAccountId"`
	Realm          string   `json:"realm"`
	Secret         string   `json:"secret"`
	SubjectID      string   `json:"subjectId"`
	Target         string   `json:"target"`
	TokenType      string   `json:"tokenType"`
}

// AzureClaims represents the claims structure for Azure JWT tokens
type AzureClaims struct {
	Aud      string   `json:"aud"`
	Email    string   `json:"email"` // User email
	Groups   []string `json:"groups,omitempty"`
	ObjectID string   `json:"oid"` // Azure Object ID (may not always be present)
}

func (apt *AzurePimToken) ComputeAdditionalFields() error {
	if apt.SubjectID != "" {
		return nil
	}

	logger := log.InitializeLogger()

	// Manually decode the JWT token to extract claims without signature verification
	// Split the JWT into parts (header.payload.signature)
	logger.WithPrefix("JWT").Info("Decoding JWT token...")
	parts := strings.Split(apt.Secret, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}
	logger.WithPrefix("JWT").Infof("JWT token found: %s***", apt.Secret[:3])
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
	apt.Audience = claims.Aud
	apt.Email = claims.Email
	apt.Groups = claims.Groups
	apt.SubjectID = claims.ObjectID
	return nil
}
