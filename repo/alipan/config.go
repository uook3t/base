package repo

type OpenConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func NewOpenConfig(clientID, clientSecret string) *OpenConfig {
	return &OpenConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}
