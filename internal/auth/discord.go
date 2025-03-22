package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"wat/internal/config"
	"wat/internal/models"
)

func GetDiscordOauthUrl(c *config.DiscordOauthConfig) string {
	redirectUri := c.RedirectUri
	oauthEndpoint := "https://discord.com/oauth2/authorize" + "?" +
		url.PathEscape(fmt.Sprintf("response_type=code&client_id=%s&scope=guilds+identify&redirect_uri=%s", c.ClientID, redirectUri))
	return oauthEndpoint
}

func GetDiscordAccessToken(code string, config *config.DiscordOauthConfig) (*models.AccessTokenDetails, error) {
	request, err := http.PostForm(
		"https://discord.com/api/v10/oauth2/token",
		url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {config.RedirectUri},
			"client_id":     {config.ClientID},
			"client_secret": {config.ClientSecret},
		},
	)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	var accessTokenDetails models.AccessTokenDetails
	err = json.Unmarshal(body, &accessTokenDetails)
	if err != nil {
		return nil, err
	}
	return &accessTokenDetails, nil
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
	User    struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Avatar        string `json:"avatar"`
		Discriminator string `json:"discriminator"`
		GlobalName    string `json:"global_name"`
		PublicFlags   int    `json:"public_flags"`
	} `json:"user"`
}

func GetDiscordAuthDetails(accessToken *models.AccessTokenDetails) (*DiscordAuthDetails, error) {
	request := http.Request{
		Method: "GET",
		URL:    &url.URL{Scheme: "https", Host: "discord.com", Path: "/api/v10/oauth2/@me"},
		Header: http.Header{
			"Authorization": {"Bearer " + accessToken.AccessToken},
			"Accept":        {"application/json"},
		},
	}
	resp, err := http.DefaultClient.Do(&request)
	if err != nil {
		fmt.Println("Unable to make request to discord", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("discord returned unauthorized", resp.StatusCode)
		return nil, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read body", err)
		return nil, err
	}
	var discordAuthDetails DiscordAuthDetails
	err = json.Unmarshal(body, &discordAuthDetails)
	if err != nil {
		fmt.Println("Unable to unmarshal discord auth details", err)
		return nil, err
	}
	return &discordAuthDetails, nil
}

type guildsResponse []struct {
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

func GetDiscordGuildMembership(accessToken *models.AccessTokenDetails, guildId string) (bool, error) {
	request := http.Request{
		Method: "GET",
		URL:    &url.URL{Scheme: "https", Host: "discord.com", Path: "/api/v10/users/@me/guilds"},
		Header: http.Header{
			"Authorization": {"Bearer " + accessToken.AccessToken},
			"Accept":        {"application/json"},
		},
	}
	resp, err := http.DefaultClient.Do(&request)
	if err != nil {
		fmt.Println("Unable to make request to discord", err)
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("discord returned unauthorized", resp.StatusCode)
		return false, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read body", err)
		return false, err
	}
	var guildsResponse guildsResponse
	err = json.Unmarshal(body, &guildsResponse)
	if err != nil {
		fmt.Println("Unable to unmarshal discord guilds", err)
		return false, err
	}
	for _, guild := range guildsResponse {
		if guildId == guild.ID {
			return true, nil
		}
	}
	return false, nil
}
