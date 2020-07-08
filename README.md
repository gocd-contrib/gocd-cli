# gocd-cli

A command-line companion to GoCD. In its current state, it mostly has helpers for developing config-repo definitions. More to come later.

## Development

This is a [golang](https://golang.org/) project, so if you are new to this, please set up your environment according to [this page](https://golang.org/doc/code.html#Workspaces).
Golang version >= 1.12.7 is needed.

To build the `gocd` binary, run the `build.sh` command:

```bash
$ ./build.sh
```

This will generate the `gocd` binary in the repository's root directory (i.e., most likely your working directory).

You will likely need to [`fetch`](#fetch-fetch-config-repo-plugins) recent versions of config-repo plugins to do anything interesting.

### Build in docker

Alternatively, if you have local docker daemon, you can build the binary using [docker golang-dojo](https://github.com/kudulab/docker-golang-dojo) image.

1. [Install docker](https://docs.docker.com/install/), if you haven't already.
2. [Install Dojo](https://github.com/ai-traders/dojo#installation), it is a self-contained binary, so just place it somewhere on the `PATH`.
On **Linux**:
```bash
DOJO_VERSION=0.4.3
wget -O dojo https://github.com/ai-traders/dojo/releases/download/${DOJO_VERSION}/dojo_linux_amd64
sudo mv dojo /usr/local/bin
sudo chmod +x /usr/local/bin/dojo
```
3. Build gocd-cli with
```
dojo ./build.sh
```
This will generate the `gocd` binary in the repository's root directory.

## Usage

The general invocation pattern for `gocd` looks like this:

```bash
gocd [global-flags] COMMAND [flags] [SUBCOMMAND] [sub-flags] [args...]
```

There are built-in help screens to the `gocd` binary if you pass in `-h` or invoke with no arguments at all.

### Global flags

#### `--config path/to/file` (equivalent short-opt `-c`)

All commands support `--config` to use a specified config file. If this file does not exist, it will be created, so long as filesystem permissions allow. This feature can be used to support multiple profiles.

If this flag is omitted, `gocd` will use `${HOME}/.gocd/settings.yaml` as the default.

For any command depending on a config value stored in a non-default config file path, one must specify the `--config` flag; `gocd` will not "remember" the config file location for subsequent invocations.

Example:

```bash
# Configure auth token in a specific file
$ gocd -c $HOME/myadminconfigfile.yaml config auth-token mysupersecrettoken

# Use the auth configuration stored in that same file to make a preflight API call to GoCD
$ gocd -c $HOME/myadminconfigfile.yaml configrepo --yaml preflight allpipelines.gocd.yaml
```

#### `--debug` (equivalent short-opt `-X`)

Enables verbose debugging output to aid in troubleshooting any issues with the `gocd` tool.

#### `--quiet` (equivalent short-opt `-q`)

Suppresses **most** output. Certain fatal errors will not be suppressed so as to not provide false negative feedback.

#### `--help` (equivalent short-opt `-h`)

Prints a help/usage message for the current command/subcommand.

### Subcommands

### `config`: Setting Configuration values for the CLI
#### Example: `server-url`: Save GoCD server base URL for API calls

Note: The base URL is the URL up through the context path, which defaults to `/go` unless configured otherwise

```bash
$ gocd config server-url https://build.gocd.org/go
```

#### Example: `auth-token`, `auth-basic`, `auth-none`: Configure auth credentials for API calls

Using API Token Authentication:

```bash
# Save API token to your config file
$ gocd config auth-token mysupersecrettoken

# Also accepts input from a shell pipe; must specify the `-` argument.
$ cat /path/to/token.txt | gocd config auth-token -

$ gocd config auth-token mysupersecrettoken
```

Using Basic Authentication:

```bash
# Save Basic Authentication credentials to your config file
$ gocd config auth-basic myuser secretpassword
```

Using No Authentication (**MUST be set when GoCD server security is disabled**):

```bash
# Configures settings to omit the `Authorization` header in API calls. This os
# mainly aimed at test GoCD server instances that have security disabled.
$ gocd config auth-none
```

#### Example: `delete`: Delete auth credentials

```bash
# Requires the config-key to delete as the only argument (auth, server-url)

# Deletes the authentication configuration
$ gocd config delete auth

# Deletes the server URL configuration
$ gocd config delete server-url
```

### Configuration by Environment Variables

Alternatively, the settings can be configured or overridden using environment variables.

#### Example: Set auth credentials with Environment Variables

```bash
# For auth-token, set GOCDCLI_AUTH.TYPE=token and GOCDCLI_AUTH.TOKEN=mysupersecrettoken
# For auth-basic, set GOCDCLI_AUTH.TYPE=basic, GOCDCLI_AUTH.USER=myuser, amd GOCDCLI_AUTH.PASSWORD=mysupersecretpasswd
# For auth-none, set GOCDCLI_AUTH.TYPE=none
# For server-url, set GOCDCLI_SERVER.URL=https://your-gocd-host/go

$ env "GOCDCLI_AUTH.TYPE=token" "GOCDCLI_AUTH.TOKEN=mysupersecrettoken" gocd configrepo --yaml preflight my-pipeline.gocd.yaml
OK
```

### `configrepo`: Pipelines as Code commands (a.k.a. "Config Repos")

#### Flags supported by all `configrepo` subcommands

* `--plugin-id` or `-i`: **REQUIRED** - Specifies the plugin ID of the config-repo plugin used to process command input.
* `--plugin-dir` or `-d`: Specifies the path containing config-repo plugins. Certain commands require a locally cached copy of the plugin jar files. Common config-repo plugins will automatically be downloaded on demand if they are not present. This defaults to `${HOME}/.gocd/plugins`
* `--yaml`: Alias for `--plugin-id yaml.config.plugin`
* `--json`: Alias for `--plugin-id json.config.plugin`
* `--groovy`: Alias for `--plugin-id cd.go.contrib.plugins.configrepo.groovy`

#### `syntax`: Syntax check
##### Example: Do a syntax check on a config-repo definition file

```bash
$ gocd configrepo --yaml syntax my-pipeline.gocd.yaml
OK
```

#### `preflight`: Preflight check
##### Example: Do a preflight check on a new config-repo definition file

```bash
# Tests the config definition as a new config-repo on the GoCD instance before committing and pushing upstream
$ gocd configrepo --yaml preflight my-pipeline.gocd.yaml
OK
```

##### Example: Do a preflight check of an existing config-repo definition file (before checking it in)

```bash
# If a config-repo already exists on your GoCD instance, you MUST specify `--repo-id YOUR_REPO_ID` (short-opt `-r`) to indicate to GoCD
# that this is testing an update to an existing configuration, and not testing a new configuration; otherwise, GoCD may report a
# duplicate pipeline name error.
$ gocd configrepo --yaml preflight -r my-existing-repo my-pipeline.gocd.yaml
OK
```

#### `fetch`: Fetch config-repo plugins

Note: the fetch command will save the plugin to `${HOME}/.gocd/plugins` or to the path specified by `--plugin-dir`

##### Example: Fetch a config-repo plugin

```bash
$ gocd configrepo --yaml fetch
Downloading https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.9.0/yaml-config-plugin-0.9.0.jar
  Fetched 3.2 MB/3.2 MB (100.0%) complete

# With a non-default `--plugin-dir`
$ gocd configrepo --yaml -d /tmp/gocd-plugins fetch
Downloading https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.9.0/yaml-config-plugin-0.9.0.jar
  Fetched 3.2 MB/3.2 MB (100.0%) complete

```

##### Example: Fetch a config-repo plugin matching a specific version or a version range

```bash
# Examples:
#
# --match-version 0.7.0 will fetch specified release 0.7.0
# --match-version '>=0.5.0 <0.8.0' # simple range; will resolve to 0.7.0
# --match-version '>=0.5.0 <0.8.0 || >=0.8.1 !0.8.3' # compound range; will resolve to 0.8.2
# --match-version '0.8.x' # wildcard, will match latest 0.8.x release

$ gocd configrepo --yaml fetch --match-version '< 0.7.0'
Downloading https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.6.2/yaml-config-plugin-0.6.2.jar
  Fetched 2.0 MB/2.0 MB (100.0%) complete
```
