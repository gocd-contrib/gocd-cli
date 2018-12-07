package utils

import (
	"archive/zip"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LocatePlugin(id string, path string) string {
	var found string

	if d, err := os.Open(path); err == nil {
		defer d.Close()

		if IsDir(path) {
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
				log.Fatal(err)
			}
		} else {
			if isPluginMatchingId(id, d.Name()) {
				found = d.Name()
			}
		}
	} else {
		log.Fatal(err)
	}

	if found == "" {
		log.Fatalf("Could not find any jars matching %s", id)
	}

	return found
}

type goplugin struct {
	Id string `xml:"id,attr"`
}

func isPluginMatchingId(id string, jar string) bool {
	r, err := zip.OpenReader(jar)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	for _, f := range r.File {
		if f.Name != "plugin.xml" {
			continue
		}

		if f.Name == "plugin.xml" {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
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
