package carbonrom

import (
	"carbonrombot/modules/utils"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"sort"
	"strings"
)

type File struct {
	Filename string
	Timestamp int64
}

type Rom struct {
	RomPath   string
	Md5Url    string
	Timestamp int64
}

func (r Rom) RomUrl() string {
	return utils.GenerateMirrorBitsUrl(r.RomPath)
}

const jsonUrl = "https://carbonrom.org/deltaindex.json"

func getInfo() (map[string]interface{}, error) {
	// Get the json
	plan, err := utils.DownloadFile(jsonUrl)
	if err != nil {
		return nil, err
	}

	//  Unmarshal it
	var devicesInterface map[string]interface{}
	err = json.Unmarshal(plan, &devicesInterface)
	//log.Print(devicesInterface)
	return devicesInterface, err
}

func getDeviceFiles(device string) ([]File, error) {
	devicesInterface, err := getInfo()
	if err != nil {
		return nil, err
	}
	// Get device files as ONE interface
	filesInterface := devicesInterface["./" + device].(interface{})
	// Decode device files interface to []Files
	var files []File
	err = mapstructure.Decode(filesInterface, &files)
	//log.Print(files[0])
	return files, err
}

func GetDeviceRoms(device string) ([]Rom, error) {
	files, err := getDeviceFiles(device)
	if err != nil {
		return nil, err
	}

	// It split .zip and .md5sum files to two Files arrays
	var md5sums []File
	var zips []File
	for _, file := range files {
		// Usually device json array consist of .zip's and .md5sum's
		// The .md5sum for the .zip is after the .zip
		// Go to split it to two arrays
		if strings.HasSuffix(file.Filename, ".zip") {
			zips = append(zips, file)
		}
		if strings.HasSuffix(file.Filename, ".md5sum") {
			md5sums = append(md5sums, file)
		}
	}
	//log.Println(zips)
	//log.Println(md5sums)

	// It's time to generate Roms!
	var Roms []Rom
	for i := 0; i < len(zips); i++ {
		Roms = append(Roms, Rom{RomPath: zips[i].Filename,
			Md5Url: utils.GenerateMirrorBitsUrl(md5sums[i].Filename),
			Timestamp: zips[i].Timestamp})
	}

	// Sort it by Timestamp!
	// First ROM is oldest
	sort.SliceStable(Roms, func(i, j int) bool {
		return Roms[i].Timestamp < Roms[j].Timestamp
	})
	// log.Println(Roms)
	return Roms, nil
}

func GetDevices() ([]string, error) {
	devicesInterface, err := getInfo()
	if err != nil {
		return nil, err
	}

	// Get devices list as keys from json
	devices := make([]string, 0, len(devicesInterface))
	for key := range devicesInterface {
		// Remove delta files
		if key == "./delta" {
			continue
		}
		// Remove "./" from key
		devices = append(devices, key[2:])
	}
	sort.Strings(devices)
	//log.Println(devices)
	return devices, nil
}