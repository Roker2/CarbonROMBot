package utils

import (
	"io/ioutil"
	"net/http"
	"os"
)

const fileHostURL = "https://mirrorbits.carbonrom.org"

// Contains tells whether a contains x.
func ContainsString(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Generate CarbonROM MirrorBits URL
func GenerateFileUrl(file string) string {
	// File is "./device/filename.extension"
	// First char of file is ".", remove it
	// FILE_HOST_URL is for situation, when you want to fast change link without recompilation
	customFileHostURL := os.Getenv("FILE_HOST_URL")
	if customFileHostURL == "" {
		return fileHostURL + file[1:]
	} else {
		return customFileHostURL + file[1:]
	}
}

// Download file from url
func DownloadFile(url string) ([]byte, error) {
	// Get the file
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//log.Println(resp.Body)

	// Read it
	return ioutil.ReadAll(resp.Body)
}