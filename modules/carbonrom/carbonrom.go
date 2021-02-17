package carbonrom

import (
	"carbonrombot/modules/utils"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"sort"
	"strings"
	"time"
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

func (r Rom) RomName() string {
	// Remove "./" from RomPath
	name := r.RomPath[2:]
	// name has structure "device/romname.zip"
	// Need to remove "device/"
	name = name[strings.Index(name, "/") + 1:]
	// Remove ".zip" and return it
	return name[:len(name) - 4]
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
	var zips []File
	for _, file := range files {
		// Usually device json array consist of .zip's and .md5sum's
		// We need to get only .zip's
		if strings.HasSuffix(file.Filename, ".zip") {
			zips = append(zips, file)
		}
	}
	//log.Println(zips)

	// It's time to generate Roms!
	var Roms []Rom
	for i := 0; i < len(zips); i++ {
		Roms = append(Roms, Rom{RomPath: zips[i].Filename,
			// Usually .md5sum files has name "romname.zip.md5sum"
			Md5Url: utils.GenerateMirrorBitsUrl(zips[i].Filename + ".md5sum"),
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