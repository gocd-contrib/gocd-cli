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
	utils.Debug(`Searching for plugin=%q in path=%q`, id, path)

	if d, err := os.Open(path); err != nil {
		return "", utils.InspectError(err, `opening plugin path %q`, path)
	} else {
		defer d.Close()

		if utils.IsDir(path) {
			utils.Debug(`path %q is a directory`, path)

			if files, err := d.Readdir(-1); err == nil {
				for _, file := range files {
					if strings.HasSuffix(file.Name(), ".jar") {
						jarFile := filepath.Join(d.Name(), file.Name())

						utils.Debug(`considering jar file %q`, jarFile)

						if found, err := isPluginMatchingId(id, jarFile); err != nil {
							return "", utils.InspectError(err, `testing jar %q for plugin id %q`, jarFile, id)
						} else {
							if found {
								utils.Debug(`Found plugin %q in jar %q`, id, jarFile)
								return jarFile, nil
							}
						}
					}
				}
			} else {
				return "", err
			}
		} else {
			utils.Debug(`path %q is a file`, path)

			if found, err := isPluginMatchingId(id, d.Name()); err != nil {
				return "", utils.InspectError(err, `testing jar %q for plugin id %q`, path, id)
			} else {
				if found {
					utils.Debug(`Found plugin %q in jar %q`, id, d.Name())
					return d.Name(), nil
				}
			}
		}

		utils.Debug(`Failed to find jar for plugin=%q in path=%q`, id, path)
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
	utils.Debug(`Testing if jar %q is plugin=%q`, jar, id)

	r, err := zip.OpenReader(jar)
	if err != nil {
		return false, utils.InspectError(err, `reading jar %q`, jar)
	}

	defer r.Close()

	for _, f := range r.File {
		if f.Name != "plugin.xml" {
			continue
		}

		if f.Name == "plugin.xml" {
			rc, err := f.Open()
			if err != nil {
				return false, utils.InspectError(err, `opening embedded plugin.xml descriptor`)
			}

			var b []byte
			b, err = ioutil.ReadAll(rc)
			if err != nil {
				return false, utils.InspectError(err, `reading embedded plugin.xml descriptor`)
			}

			var pl goplugin
			if err = xml.Unmarshal(b, &pl); err != nil {
				return false, utils.InspectError(err, "parsing embedded plugin.xml descriptor:\n%s", string(b))
			}

			if pl.Id == id {
				utils.Debug(`jar %q matches plugin %q`, jar, id)
				return true, nil
			}
		}
	}

	utils.Debug(`jar %q does not match plugin %q`, jar, id)
	return false, nil
}
