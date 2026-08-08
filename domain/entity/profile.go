package entity

import (
	"encoding/json"
	"time"
)

type Profile struct {
	ID                  int             `json:"id"`
	Name                string          `json:"name"`
	Slogan              string          `json:"slogan"`
	ShortDescription    json.RawMessage `json:"short_description"`
	DetailedDescription json.RawMessage `json:"detailed_description"`
	Address             string          `json:"address"`
	Role                string          `json:"role"`
	Status              string          `json:"status"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type SocialMedia struct {
	ID           int    `json:"id"`
	PlatformName string `json:"platform_name"`
	IconURL      string `json:"icon_url"`
	LinkURL      string `json:"link_url"`
	SortOrder    int    `json:"sort_order"`
}