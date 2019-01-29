# gocd-cli

A command-line companion to GoCD. In its current state, it mostly has helpers for developing config-repo definitions. More to come later.

## Development

This is a [golang](https://golang.org/) project, so if you are new to this, please set up your environment according to [this page](https://golang.org/doc/code.html#Workspaces).

To build the `gocd` binary, run the `build.sh` command:

```bash
$ ./build.sh
```

This will generate the `gocd` binary in the repository's root directory (i.e., most likely your working directory).

You will likely need to [`fetch`](#Fetch-plugins) recent versions of config-repo plugins to do anything interesting.

## Usage

There are built-in help screens to the `gocd` binary if you pass in `-h` or invoke with no arguments at all.

### Syntax check
#### Example: Do a syntax check on a config-repo definition file

```
# required flag plugin id (--yaml or --json or --groovy alias or -i 'other.plugin.id')

$ ./gocd configrepo -i yaml.config.plugin syntax my-pipeline.gocd.yaml
OK
```

### Preflight check
#### Example: Do a preflight check on a new config-repo definition file

```
# required flag plugin id (--yaml or --json or --groovy alias or -i 'other.plugin.id')
# required argument(s) filename(s)

$ ./gocd configrepo --yaml preflight my-pipeline.gocd.yaml
OK
```

#### Example: Do a preflight check of an existing config-repo definition file (before checking it in)

```
# required flag plugin id (--yaml or --json or --groovy alias or -i 'other.plugin.id')
# optional flag --repo-id 'repo-id' or -r 'repo-id'
# required argument(s) filename(s)

$ ./gocd configrepo --yaml -r 'my-existing-repo' preflight my-pipeline.gocd.yaml
OK
```

### Fetch plugins

Note: the fetch command will save the plugin to $HOME/.gocd/plugins

#### Example: Fetch a config-repo plugin

```
# required flag plugin id (--yaml or --json or --groovy alias or -i 'other.plugin.id')
# optional flag --match-version (omitting will fetch the latest)

$ ./gocd configrepo -i yaml.config.plugin fetch
Downloading https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.8.3/yaml-config-plugin-0.8.3.jar
  Fetched 2.3 MB/2.3 MB (100.0%) complete
```

#### Example: Fetch a config-repo plugin matching a specific version or a version range

```
# required flag plugin id (--yaml or --json or --groovy alias or -i 'other.plugin.id')
# optional flag --match-version accepts a semver (semantic version) query or range string
#
# Examples:
#
# --match-version 0.7.0 will fetch specified release 0.7.0
# --match-version '>=0.5.0 <0.8.0' # simple range; will resolve to 0.7.0
# --match-version '>=0.5.0 <0.8.0 || >=0.8.1 !0.8.3' # compound range; will resolve to 0.8.2
# --match-version '0.8.x' # wildcard, will match latest 0.8.x release

$ ./gocd configrepo -i yaml.config.plugin fetch --match-version '< 0.7.0'
Downloading https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.6.2/yaml-config-plugin-0.6.2.jar
  Fetched 2.0 MB/2.0 MB (100.0%) complete
```

## Configuration

### Setting Configuration with CLI
#### Example: Save GoCD server base URL for API calls

Note: The base URL is the URL up through the context path, which defaults to `/go` unless configured otherwise

```
# required argument url
# optional flag --config 'path/to/file' or -c 'path/to/file' (omitting will use default file $HOME/.gocd/settings.yaml)

$ ./gocd config server-url https://build.gocd.org/go
```

#### Example: Save auth credentials (currently, basic auth only) for API calls

```
# required arguments username password
# optional flag --config 'path/to/file' or -c 'path/to/file'
# note: if the -config flag is used, then all subsequent commands
# dependent on that configuration must also include the flag

$ ./gocd config -c '${HOME}/.gocd/myadminconfigfile.yaml auth-basic myuser secretpassword
```

#### Example: Delete auth credentials

```
# required argument config-key to delete (auth, server-url)
# optional flag --config 'path/to/file' or -c 'path/to/file' (omitting will use default file $HOME/.gocd/settings.yaml)

$ ./gocd config delete auth
```

### Using Environment Variables

Alternatively the settings can be configured or overridden using environment variables.

#### Example: Set auth credentials with Environment Variables

```
# required variables GOCDCLI_AUTH.TYPE, GOCDCLI_AUTH.USER, GOCDCLI_AUTH.PASSWORD

$ env "GOCDCLI_AUTH.TYPE=basic" env "GOCDCLI_AUTH.USER=myuser" env "GOCDCLI_AUTH.PASSWORD=mypassword" ./gocd configrepo --yaml preflight my-pipeline.gocd.yaml
OK
```