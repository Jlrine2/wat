package models

import "time"

type Session struct {
	AccessToken  string              `json:"access_token"`
	ExpiresIn    int                 `json:"expires_in"`
	RefreshToken string              `json:"refresh_token"`
	Scope        string              `json:"scope"`
	User         *DiscordAuthDetails `json:"user"`
	Guilds       []UserGuild         `json:"guilds"`
}

type DiscordAuthDetails struct {
	Application struct {
		ID                  string `json:"id"`
		Name                string `json:"name"`
		Icon                string `json:"icon"`
		Description         string `json:"description"`
		Hook                bool   `json:"hook"`
		BotPublic           bool   `json:"bot_public"`
		BotRequireCodeGrant bool   `json:"bot_require_code_grant"`
		VerifyKey           string `json:"verify_key"`
	} `json:"application"`
	Scopes  []string  `json:"scopes"`
	Expires time.Time `json:"expires"`
	User    *struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Avatar        string `json:"avatar"`
		Discriminator string `json:"discriminator"`
		GlobalName    string `json:"global_name"`
		PublicFlags   int    `json:"public_flags"`
	} `json:"user"`
}

type UserGuild struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Icon                     string   `json:"icon"`
	Banner                   string   `json:"banner"`
	Owner                    bool     `json:"owner"`
	Permissions              string   `json:"permissions"`
	Features                 []string `json:"features"`
	ApproximateMemberCount   int      `json:"approximate_member_count"`
	ApproximatePresenceCount int      `json:"approximate_presence_count"`
}
