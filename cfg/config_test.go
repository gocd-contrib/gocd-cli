package cfg

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const TEST_CONF_FILE = CONFIG_FILENAME + "." + CONFIG_FILETYPE
const TEST_URL = "http://localhost:1234/go/test"
const TEST_USER = "admin"
const TEST_PASSWORD = "badger"

type keyVal map[string]interface{}

func writeContent(fs afero.Fs, path, content string) error {
	f, e := fs.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))

	if e != nil {
		return e
	}

	defer f.Close()
	_, e = io.Copy(f, strings.NewReader(content))
	return e
}

// creates a test Config instance, backed by a memory mapped afero.Fs
// if create == true, creates an empty config file on the afero.Fs
func testConf(create bool) *Config {
	v := viper.New()
	v.SetEnvPrefix(CONFIG_ENV_PFX)
	v.AutomaticEnv()
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	v.SetConfigType(CONFIG_FILETYPE)
	v.SetConfigName(CONFIG_FILENAME)
	v.SetConfigFile(TEST_CONF_FILE)

	if create {
		v.WriteConfigAs(TEST_CONF_FILE)
	}

	return &Config{native: v, fs: fs}
}

func TestGetBasicAuth(t *testing.T) {
	as := asserts(t)
	c := testConf(true)
	c.native.Set("auth.type", "basic")
	c.native.Set("auth.user", TEST_USER)
	c.native.Set("auth.password", TEST_PASSWORD)
	c.native.WriteConfig()

	as.deepEq(map[string]string{
		"type":     "basic",
		"user":     TEST_USER,
		"password": TEST_PASSWORD,
	}, c.GetAuth())
}

func TestSetBasicAuth(t *testing.T) {
	as := asserts(t)
	c := testConf(true)

	as.configEq(make(keyVal, 0), c.fs)

	as.ok(c.SetBasicAuth(TEST_USER, TEST_PASSWORD))

	as.configEq(keyVal{
		"auth": map[string]string{
			"type":     "basic",
			"user":     TEST_USER,
			"password": TEST_PASSWORD,
		},
	}, c.fs)
}

func TestSetBasicAuthShouldValidatePrescenseOfUserAndPassword(t *testing.T) {
	as := asserts(t)
	c := testConf(true)

	as.err("Must specify user and password", c.SetBasicAuth("", TEST_PASSWORD))
	as.err("Must specify user and password", c.SetBasicAuth(TEST_USER, ""))
	as.err("Must specify user and password", c.SetBasicAuth("", ""))
	as.eq(0, len(c.GetAuth()))
}

func TestSetServerURL(t *testing.T) {
	as := asserts(t)
	c := testConf(true)

	as.configEq(make(keyVal, 0), c.fs)

	as.ok(c.SetServerUrl(TEST_URL))

	as.configEq(keyVal{
		"server": map[string]string{
			"url": TEST_URL,
		},
	}, c.fs)
}

func TestGetServerURL(t *testing.T) {
	as := asserts(t)
	c := testConf(true)
	c.native.Set("server.url", TEST_URL)
	c.native.WriteConfig()

	as.eq(TEST_URL, c.GetServerUrl())
}

func TestSetServerURLValidatesURL(t *testing.T) {
	as := asserts(t)
	c := testConf(true)

	as.err("Must specify a url", c.SetServerUrl(""))
	as.err("URL must include protocol and hostname", c.SetServerUrl("foo.bar"))
	as.err("URL must include protocol and hostname", c.SetServerUrl("http://"))
	as.err("Port must be numeric", c.SetServerUrl("http://localhost:foo/bar"))
	as.eq("", c.GetServerUrl())
}

func TestConsumeCreatesDefaultConfigFileIfNotExist(t *testing.T) {
	as := asserts(t)
	c := testConf(false)

	baseDir, e := homedir.Dir()
	as.ok(e)

	expectedPath := filepath.Join(baseDir, CONFIG_DIRNAME, TEST_CONF_FILE)
	exists, err := afero.Exists(c.fs, expectedPath)

	as.ok(err)
	as.not(exists)

	as.ok(c.Consume(""))

	exists, err = afero.Exists(c.fs, expectedPath)
	as.ok(err)
	as.is(exists)

	isDir, er := afero.IsDir(c.fs, expectedPath)
	as.ok(er)
	as.not(isDir)
}

func TestConsumeReturnsErrorIfSpecifiedConfigFileNotExists(t *testing.T) {
	as := asserts(t)
	c := testConf(false)

	expectedPath := "foo.yaml"
	exists, err := afero.Exists(c.fs, expectedPath)

	as.ok(err)
	as.not(exists)

	as.err("open "+expectedPath+": file does not exist", c.Consume(expectedPath))

	exists, err = afero.Exists(c.fs, expectedPath)
	as.ok(err)
	as.not(exists)
}

func TestConsumeLoadsConfigFileAtDefaultPathWhenExists(t *testing.T) {
	as := asserts(t)
	c := testConf(false)

	baseDir, e := homedir.Dir()
	as.ok(e)

	expectedPath := filepath.Join(baseDir, CONFIG_DIRNAME, TEST_CONF_FILE)

	as.ok(writeContent(c.fs, expectedPath, `server:
  url: http://test-server:8080
`))

	exists, err := afero.Exists(c.fs, expectedPath)
	as.ok(err)
	as.is(exists)

	as.ok(c.Consume(""))
	as.eq("http://test-server:8080", c.GetServerUrl())
}

func TestConsumeLoadsSpecifiedConfigFileWhenExists(t *testing.T) {
	as := asserts(t)
	c := testConf(false)

	expectedPath := "/etc/foo.yaml"

	as.ok(writeContent(c.fs, expectedPath, `server:
  url: http://test-server:8080
`))

	exists, err := afero.Exists(c.fs, expectedPath)
	as.ok(err)
	as.is(exists)

	as.ok(c.Consume(expectedPath))
	as.eq("http://test-server:8080", c.GetServerUrl())
}
