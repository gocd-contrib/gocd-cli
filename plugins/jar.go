package plugins

import (
	"archive/zip"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
)

func PluginById(id string, path string) (found string) {
	if d, err := os.Open(path); err == nil {
		defer d.Close()

		if utils.IsDir(path) {
			if files, err := d.Readdir(-1); err == nil {
				for _, file := range files {
					if strings.HasSuffix(file.Name(), ".jar") {
						jarFile := filepath.Join(d.Name(), file.Name())

						if isPluginMatchingId(id, jarFile) {
							found = jarFile
							break
						}
					}
				}
			} else {
				utils.AbortLoudly(err)
			}
		} else {
			if isPluginMatchingId(id, d.Name()) {
				found = d.Name()
			}
		}
	} else {
		utils.AbortLoudly(err)
	}
	return
}

func LocatePlugin(id string, path string) string {
	found := PluginById(id, path)

	if found == "" {
		utils.DieLoudly(1, "Could not find any plugin jars with id: %s", id)
	}

	return found
}

type goplugin struct {
	Id    string `xml:"id,attr"`
	About about  `xml:"about"`
}

type about struct {
	XMLName xml.Name `xml:"about"`
	Name    string   `xml:"name"`
	Version string   `xml:"version"`
}

func isPluginMatchingId(id string, jar string) bool {
	r, err := zip.OpenReader(jar)
	if err != nil {
		utils.AbortLoudly(err)
	}

	defer r.Close()

	for _, f := range r.File {
		if f.Name != "plugin.xml" {
			continue
		}

		if f.Name == "plugin.xml" {
			rc, err := f.Open()
			if err != nil {
				utils.AbortLoudly(err)
			}

			b, _ := ioutil.ReadAll(rc)

			var pl goplugin
			xml.Unmarshal(b, &pl)

			if pl.Id == id {
				return true
			}
		}
	}

	return false
}
