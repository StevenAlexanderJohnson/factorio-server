package models

// AdminList represents a list of admin player usernames.
type AdminList []string

// Whitelist represents a list of whitelisted player usernames.
type Whitelist []string

// BanEntry represents a banned user entry in the Factorio banlist.
type BanEntry struct {
	Username string `json:"username"`
	Reason   string `json:"reason,omitempty"`
}

// BanList represents a list of banned users.
type BanList []string
