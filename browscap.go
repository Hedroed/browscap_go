package browscap_go

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	DownloadUrl     = "http://browscap.org/stream?q=PHP_BrowsCapINI"
	CheckVersionUrl = "http://browscap.org/version-number"
)

var (
	dict        *dictionary
	initialized bool
	version     string
	debug       bool
)

func Debug(val bool) {
	debug = val
}

func InitBrowsCap(path string, force bool) error {
	if initialized && !force {
		return nil
	}
	var err error

	// Load ini file
	if dict, err = loadFromIniFile(path); err != nil {
		return fmt.Errorf("browscap: An error occurred while reading file, %v ", err)
	}

	if verDictionary, exists := dict.mapped["GJK_Browscap_Version"]; exists {
		version = verDictionary["Version"]
	}

	initialized = true
	return nil
}

func InitializedVersion() string {
	return version
}

func LastVersion() (string, error) {
	response, err := http.Get(CheckVersionUrl)
	if err != nil {
		return "", fmt.Errorf("browscap: error sending request, %v", err)
	}
	defer response.Body.Close()

	// Get body of response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("browscap: error reading the response data of request, %v", err)
	}

	// Check 200 status
	if response.StatusCode != 200 {
		return "", fmt.Errorf("browscap: error unexpected status code %d", response.StatusCode)
	}

	return string(body), nil
}

func DownloadFile(saveAs string) error {
	response, err := http.Get(DownloadUrl)
	if err != nil {
		return fmt.Errorf("browscap: error sending request, %v", err)
	}
	defer response.Body.Close()

	// Get body of response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("browscap: error reading the response data of request, %v", err)
	}

	// Check 200 status
	if response.StatusCode != 200 {
		return fmt.Errorf("browscap: error unexpected status code %d", response.StatusCode)
	}

	if err = ioutil.WriteFile(saveAs, body, os.ModePerm); err != nil {
		return fmt.Errorf("browscap: error saving file, %v", err)
	}

	return nil
}

func GetBrowserData(userAgent string) (map[string]string, bool) {
	return searchIndexedBrowserData(userAgent)
}

func searchIndexedBrowserData(userAgent string) (map[string]string, bool) {
	agent := strings.ToLower(userAgent)

	best := dict.tree.Find(agent)

	if best == "" {
		return nil, false
	}

	data := dict.getData(best)
	return data, true
}
