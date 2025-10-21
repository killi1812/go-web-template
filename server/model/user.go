package model

// User represents a user in the lobby.
// Based on https://discord.com/developers/docs/resources/user#user-object
type User struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	GlobalName    *string `json:"global_name,omitempty"`
	Avatar        *string `json:"avatar,omitempty"`
	Bot           *bool   `json:"bot,omitempty"`
	System        *bool   `json:"system,omitempty"`
	MFAEnabled    *bool   `json:"mfa_enabled,omitempty"`
	Banner        *string `json:"banner,omitempty"`
	AccentColor   *int    `json:"accent_color,omitempty"`
	Locale        *string `json:"locale,omitempty"`
	Verified      *bool   `json:"verified,omitempty"`
	Email         *string `json:"email,omitempty"`
	Flags         *int    `json:"flags,omitempty"`
	PremiumType   *int    `json:"premium_type,omitempty"`
	PublicFlags   *int    `json:"public_flags,omitempty"`
	// NOTE:
	// Complex fields like avatar_decoration_data, collectibles, and primary_guild are omitted for simplicity.
	// They can be added as their own structs if needed.
}
