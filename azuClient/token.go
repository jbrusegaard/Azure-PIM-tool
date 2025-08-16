package azuClient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
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
}

// AzureClaims represents the claims structure for Azure JWT tokens
type AzureClaims struct {
	Subject    string `json:"sub"`         // Standard JWT subject identifier
	ObjectID   string `json:"oid"`         // Azure Object ID (may not always be present)
	Email      string `json:"email"`       // User email
	UniqueName string `json:"unique_name"` // Unique user identifier
	TenantID   string `json:"tid"`         // Azure Tenant ID
	jwt.StandardClaims
}

func (apt *AzurePimToken) ComputeSubjectID() error {
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

	// Debug: Print the raw payload to see what's in the JWT
	fmt.Println("Raw JWT payload:", string(decodedPayload))

	// Parse the JSON claims into a generic map first to see all fields
	var rawClaims map[string]interface{}
	if err := json.Unmarshal(decodedPayload, &rawClaims); err != nil {
		return fmt.Errorf("failed to unmarshal JWT claims: %v", err)
	}

	// Debug: Print all available claims
	fmt.Println("Available claims:")
	for key, value := range rawClaims {
		fmt.Printf("  %s: %v\n", key, value)
	}

	// Parse into our custom claims structure
	var claims AzureClaims
	if err := json.Unmarshal(decodedPayload, &claims); err != nil {
		return fmt.Errorf("failed to unmarshal JWT claims: %v", err)
	}

	// Debug: Check what we got in our claims struct
	fmt.Printf("AzureClaims - Subject: '%s', ObjectID: '%s', Email: '%s', UniqueName: '%s'\n",
		claims.Subject, claims.ObjectID, claims.Email, claims.UniqueName)

	// Use the Subject field (sub) as the primary SubjectID
	if claims.Subject != "" {
		apt.SubjectID = claims.Subject
		fmt.Println("Using 'sub' field as SubjectID:", claims.Subject)
	} else if claims.ObjectID != "" {
		// Fallback to ObjectID if available
		apt.SubjectID = claims.ObjectID
		fmt.Println("Using 'oid' field as SubjectID:", claims.ObjectID)
	} else if claims.Email != "" {
		// Fallback to email
		apt.SubjectID = claims.Email
		fmt.Println("Using 'email' field as SubjectID:", claims.Email)
	} else if claims.UniqueName != "" {
		// Fallback to unique name
		apt.SubjectID = claims.UniqueName
		fmt.Println("Using 'unique_name' field as SubjectID:", claims.UniqueName)
	} else {
		return fmt.Errorf("no suitable subject identifier found in JWT claims")
	}

	fmt.Println("Final SubjectID:", apt.SubjectID)
	return nil
}
