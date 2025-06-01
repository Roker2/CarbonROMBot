package carbonrom

import (
	"carbonrombot/modules/utils"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

type File struct {
	Filename  string
	Timestamp int64
}

type Rom struct {
	RomPath   string
	Md5Url    string
	Timestamp int64
}

var Versions = map[string]string{
	"CARBON-CR-6.1":  "Android 8.1 (Oreo)",
	"CARBON-CR-7.0":  "Android 9 (Pie)",
	"CARBON-CR-8.0":  "Android 10 (Q)",
	"CARBON-CR-9.0":  "Android 11 (R)",
	"CARBON-CR-11.0": "Android 13 (T)",
}

func (r Rom) RomUrl() string {
	return utils.GenerateFileUrl(r.RomPath)
}

func (r Rom) RomName() string {
	// Remove "./" from RomPath
	name := r.RomPath[2:]
	// name has structure "device/romname.zip"
	// Need to remove "device/"
	name = name[strings.Index(name, "/")+1:]
	// Remove ".zip" and return it
	return name[:len(name)-4]
}

func (r Rom) RomVersion() (string, error) {
	for carbonVersion, androidVersion := range Versions {
		if strings.HasPrefix(r.RomName(), carbonVersion) {
			return androidVersion, nil
		}
	}
	// If version is not in the map, return "nullptr" (easter egg) and error
	return "nullptr", fmt.Errorf("carbonrom error: I can not find Android version for %s", r.RomName())
}

func (r Rom) Md5() (string, error) {
	// Get MD5 from url
	file, err := utils.DownloadFile(r.Md5Url)
	if err != nil {
		return "", err
	}
	// .md5sum file has structure "md5sum OUT/ROM.zip"
	// Second part is unneeded, remove it via splitting
	return strings.Split(string(file), " ")[0], err
}

func (r Rom) GetDateAsString() string {
	return time.Unix(r.Timestamp, 0).String()
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
	return devicesInterface, err
}

func getDeviceFiles(device string) ([]File, error) {
	devicesInterface, err := getInfo()
	if err != nil {
		return nil, err
	}
	// Get device files as ONE interface
	filesInterface := devicesInterface["./"+device].(interface{})
	// Decode device files interface to []Files
	var files []File
	err = mapstructure.Decode(filesInterface, &files)
	return files, err
}

func GetDeviceRoms(device string) ([]Rom, error) {
	files, err := getDeviceFiles(device)
	if err != nil {
		return nil, err
	}

	// It split .zip and .md5sum files to two Files arrays
	var zips []File
	for _, file := range files {
		// Usually device json array consists of .zip's and .md5sum's
		// We need to get only .zip's
		if strings.HasSuffix(file.Filename, ".zip") {
			zips = append(zips, file)
		}
	}

	// It's time to generate Roms!
	var Roms []Rom
	for i := 0; i < len(zips); i++ {
		Roms = append(Roms, Rom{RomPath: zips[i].Filename,
			// Usually .md5sum files has name "romname.zip.md5sum"
			Md5Url:    utils.GenerateFileUrl(zips[i].Filename + ".md5sum"),
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
		if key == "./delta" || key == "0" {
			continue
		}
		// Remove "./" from key
		devices = append(devices, key[2:])
	}
	sort.Strings(devices)
	return devices, nil
}
