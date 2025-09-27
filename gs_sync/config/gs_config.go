package config

import (
	"encoding/json"
	"os"
)

type GSConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

// LoadFromFile loads the config from a JSON file
func (c *GSConfig) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

func BuildFromFile(filename string) *GSConfig {
	cfg := &GSConfig{}
	err := cfg.LoadFromFile(filename)
	if err != nil {
		// If loading from file fails, try loading from environment variables
		cfg.LoadFromEnv()
	}
	return cfg
}

// LoadFromEnv loads the config from environment variables
func (c *GSConfig) LoadFromEnv() {
	if val := os.Getenv("TYPE"); val != "" {
		c.Type = val
	}
	if val := os.Getenv("PROJECT_ID"); val != "" {
		c.ProjectID = val
	}
	if val := os.Getenv("PRIVATE_KEY_ID"); val != "" {
		c.PrivateKeyID = val
	}
	if val := os.Getenv("PRIVATE_KEY"); val != "" {
		c.PrivateKey = val
	}
	if val := os.Getenv("CLIENT_EMAIL"); val != "" {
		c.ClientEmail = val
	}
	if val := os.Getenv("CLIENT_ID"); val != "" {
		c.ClientID = val
	}
	if val := os.Getenv("AUTH_URI"); val != "" {
		c.AuthURI = val
	}
	if val := os.Getenv("TOKEN_URI"); val != "" {
		c.TokenURI = val
	}
	if val := os.Getenv("AUTH_PROVIDER_X509_CERT_URL"); val != "" {
		c.AuthProviderX509CertURL = val
	}
	if val := os.Getenv("CLIENT_X509_CERT_URL"); val != "" {
		c.ClientX509CertURL = val
	}
	if val := os.Getenv("UNIVERSE_DOMAIN"); val != "" {
		c.UniverseDomain = val
	}
}

// ToJSON returns the config as a JSON byte slice
func (c *GSConfig) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}