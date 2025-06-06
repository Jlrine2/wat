package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port           string `json:"port" yaml:"port"`
	Host           string `json:"host" yaml:"host"`
	MediaLocation  string `json:"media_location" yaml:"media_location"`
	ClientLocation string `json:"client_location" yaml:"client_location"`
	AuthEnabled    bool   `json:"auth_enabled" yaml:"auth_enabled"`
}

type DiscordOauthConfig struct {
	ClientID       string   `json:"client_id" yaml:"client_id"`
	ClientSecret   string   `json:"client_secret" yaml:"client_secret"`
	RedirectUri    string   `json:"redirect_uri" yaml:"redirect_uri"`
	MemberGuildIds []string `json:"member_guild_ids" yaml:"member_guild_ids"`
	AdminUserIds   []string `json:"admin_user_ids" yaml:"admin_user_ids"`
}

type Config struct {
	Server       *ServerConfig       `json:"server" yaml:"server"`
	DiscordOauth *DiscordOauthConfig `json:"discord" yaml:"discord"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
