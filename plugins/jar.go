package plugins

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
)

type PluginNotFoundError struct {
	PluginId, Path string
}

func (e *PluginNotFoundError) Error() string {
	return fmt.Sprintf(`No matching plugin jar with id %q found in path %q`, e.PluginId, e.Path)
}

func PluginById(id string, path string) (string, error) {
	if d, err := os.Open(path); err != nil {
		return "", err
	} else {
		defer d.Close()

		if utils.IsDir(path) {
			if files, err := d.Readdir(-1); err == nil {
				for _, file := range files {
					if strings.HasSuffix(file.Name(), ".jar") {
						jarFile := filepath.Join(d.Name(), file.Name())

						if found, err := isPluginMatchingId(id, jarFile); err != nil {
							return "", err
						} else {
							if found {
								return jarFile, nil
							}
						}
					}
				}
			} else {
				return "", err
			}
		} else {
			if found, err := isPluginMatchingId(id, d.Name()); err != nil {
				return "", err
			} else {
				if found {
					return d.Name(), nil
				}
			}
		}
		return "", &PluginNotFoundError{Path: d.Name(), PluginId: id}
	}
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

func isPluginMatchingId(id string, jar string) (bool, error) {
	r, err := zip.OpenReader(jar)
	if err != nil {
		return false, err
	}

	defer r.Close()

	for _, f := range r.File {
		if f.Name != "plugin.xml" {
			continue
		}

		if f.Name == "plugin.xml" {
			rc, err := f.Open()
			if err != nil {
				return false, err
			}

			var b []byte
			b, err = ioutil.ReadAll(rc)
			if err != nil {
				return false, err
			}

			var pl goplugin
			if err = xml.Unmarshal(b, &pl); err != nil {
				return false, err
			}

			if pl.Id == id {
				return true, nil
			}
		}
	}

	return false, nil
}
