package utils

const mirrorbits = "https://mirrorbits.carbonrom.org/"

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