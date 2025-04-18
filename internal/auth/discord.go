package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"wat/internal/config"
	"wat/internal/models"
)

func GetDiscordOauthUrl(c *config.DiscordOauthConfig) string {
	redirectUri := c.RedirectUri
	oauthEndpoint := "https://discord.com/oauth2/authorize" + "?" +
		url.PathEscape(fmt.Sprintf("response_type=code&client_id=%s&scope=guilds+identify&redirect_uri=%s", c.ClientID, redirectUri))
	return oauthEndpoint
}

func GetDiscordAccessToken(code string, config *config.DiscordOauthConfig) (*models.Session, error) {
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
	var accessTokenDetails models.Session
	err = json.Unmarshal(body, &accessTokenDetails)
	if err != nil {
		return nil, err
	}
	return &accessTokenDetails, nil
}

func GetDiscordAuthDetails(accessToken *models.Session) (*models.DiscordAuthDetails, error) {
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
		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Println("Rate limit headers:")
			for key, values := range resp.Header {
				if strings.HasPrefix(key, "X-RateLimit") {
					fmt.Printf("%s: %s\n", key, strings.Join(values, ", "))
				}
			}

			return nil, errors.New("rate limited")
		}
		fmt.Println("discord returned unauthorized", resp.StatusCode)
		return nil, errors.New("discord returned unauthorized")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read body", err)
		return nil, err
	}
	var discordAuthDetails models.DiscordAuthDetails
	err = json.Unmarshal(body, &discordAuthDetails)
	if err != nil {
		fmt.Println("Unable to unmarshal discord auth details", err)
		return nil, err
	}
	return &discordAuthDetails, nil
}

func GetDiscordGuildMembership(accessToken *models.Session) ([]models.UserGuild, error) {
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
	var guildsResponse []models.UserGuild
	err = json.Unmarshal(body, &guildsResponse)
	if err != nil {
		fmt.Println("Unable to unmarshal discord guilds", err)
		return nil, err
	}
	return guildsResponse, nil
}
