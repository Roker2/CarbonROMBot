package utils

import (
	"io/ioutil"
	"net/http"
)

const mirrorbits = "https://hosting.carbonrom.org"

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
func GenerateMirrorBitsUrl(file string) string {
	// File is "./device/filename.extension"
	// First char of file is ".", remove it
	return mirrorbits + file[1:]
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
	plan, err := ioutil.ReadAll(resp.Body)
	return plan, err
}