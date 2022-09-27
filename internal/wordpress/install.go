package wordpress

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type PluginEntry struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Update  string `json:"update"`
	Version string `json:"Version"`
}
type SnapshotResponse struct {
	Hashes  map[string]string `json:"hashes"`
	Plugins []PluginEntry     `json:"plugins"`
}

type InstallSyncResponse struct {
	Plugins     []PluginEntry `json:"plugins"`
	SuperAdmins []string      `json:"superAdmins"`
	Users       []UserAccount `json:"users"`
	IsMultisite bool          `json:"isMultisite"`
	PrimaryURL  string        `json:"primaryURL"`
	Installs    []Install     `json:"installs"`
}

type UserAccount struct {
	ID       string `json:"user_id"`
	Nicename string `json:"display_name"`
	Username string `json:"user_login"`
	Email    string `json:"user_email"`
	Role     string `json:"roles"`
}

type Install struct {
	Users   []UserAccount `json:"users"`
	Plugins []PluginEntry `json:"plugins"`
	SiteURL string        `json:"siteURL"`
}

func GetUsersForSite(site string) (accounts []UserAccount, err error) {
	out, err := exec.Command("wp", "user", "list", "--format=json", "--url="+site, "--fields=ID,user_email,user_login,display_name,roles", "--path=public").Output()
	if err != nil {
		return accounts, err
	}
	err = json.Unmarshal(out, &accounts)
	if err != nil {
		return accounts, err
	}
	return accounts, nil
}

func GetUsers() (accounts []UserAccount, err error) {
	out, err := exec.Command("wp", "user", "list", "--format=json", "--fields=ID,user_email,user_login,display_name,roles", "--path=public").Output()
	if err != nil {
		return accounts, err
	}
	err = json.Unmarshal(out, &accounts)
	if err != nil {
		return accounts, err
	}
	return accounts, nil
}

type Site struct {
	BlogID string `json:"blog_id"`
	Url    string `json:"url"`
}

func GetSuperAdmins() ([]string, error) {
	out, err := exec.Command("wp", "super-admin", "list", "--path=public").Output()
	if err != nil {
		return []string{}, err
	}
	parts := strings.Split(string(out), "\n")
	return parts, err
}

func GetSites() ([]Site, error) {
	out, err := exec.Command("wp", "site", "list", "--format=json", "--path=public").Output()
	if err != nil {
		return []Site{}, err
	}
	entries := []Site{}
	err = json.Unmarshal(out, &entries)
	if err != nil {
		return []Site{}, err
	}
	return entries, nil
}

func GetPluginsForSite(url string) ([]PluginEntry, error) {
	out, err := exec.Command("wp", "plugin", "list", "--format=json", "--url="+url, "--path=public").Output()
	if err != nil {
		return []PluginEntry{}, err
	}
	entries := []PluginEntry{}
	err = json.Unmarshal(out, &entries)
	if err != nil {
		return []PluginEntry{}, err
	}
	return entries, nil
}

func GetPlugins() ([]PluginEntry, error) {
	out, err := exec.Command("wp", "plugin", "list", "--format=json", "--path=public").Output()
	if err != nil {
		return []PluginEntry{}, err
	}
	entries := []PluginEntry{}
	err = json.Unmarshal(out, &entries)
	if err != nil {
		return []PluginEntry{}, err
	}
	return entries, nil
}

func GetPrimaryURL() (string, error) {
	out, err := exec.Command("wp", "option", "get", "siteurl", "--path=public").Output()
	return string(out), err
}

func IsMultisite() bool {
	out, err := exec.Command("wp", "config", "get", "WP_ALLOW_MULTISITE", "--path=public").Output()
	logrus.WithField("out", out).Info("is multisite")
	return (err == nil)
}

type WPConfig struct {
	ProtectedPaths []string
}

func GenerateResponse(config WPConfig) (InstallSyncResponse, error) {
	isMulti := IsMultisite()
	if isMulti {
		return getInfoMultisite()
	}
	return getInfoSinglesite()
}

func getInfoSinglesite() (InstallSyncResponse, error) {
	users, err := GetUsers()
	if err != nil {
		return InstallSyncResponse{}, err
	}
	primaryURL, err := GetPrimaryURL()
	if err != nil {
		return InstallSyncResponse{}, err
	}
	plugins, err := GetPlugins()
	if err != nil {
		return InstallSyncResponse{}, err
	}
	return InstallSyncResponse{
		Plugins:     plugins,
		Users:       users,
		PrimaryURL:  primaryURL,
		IsMultisite: false,
	}, nil
}

func getInfoMultisite() (InstallSyncResponse, error) {
	users, err := GetUsers()
	if err != nil {
		return InstallSyncResponse{}, err
	}
	primaryURL, err := GetPrimaryURL()
	if err != nil {
		return InstallSyncResponse{}, err
	}
	plugins, err := GetPlugins()
	if err != nil {
		return InstallSyncResponse{}, err
	}
	sites, err := GetSites()
	if err != nil {
		return InstallSyncResponse{}, err
	}
	supers, err := GetSuperAdmins()
	if err != nil {
		return InstallSyncResponse{}, err
	}

	installs := []Install{}
	for _, site := range sites {
		if site.BlogID != "0" {
			p, err := GetPluginsForSite(site.Url)
			if err != nil {
				return InstallSyncResponse{}, err
			}
			u, err := GetUsersForSite(site.Url)
			if err != nil {
				return InstallSyncResponse{}, err
			}
			installs = append(installs, Install{Plugins: p, Users: u, SiteURL: site.Url})
		}
	}
	return InstallSyncResponse{
		Plugins:     plugins,
		Users:       users,
		SuperAdmins: supers,
		PrimaryURL:  primaryURL,
		Installs:    installs,
		IsMultisite: true,
	}, nil
}
