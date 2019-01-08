# gocd-cli

A command-line companion to GoCD. In its current state, it mostly has helpers for developing config-repo definitions. More to come later.

## Development

This is a [golang](https://golang.org/) project, so if you are new to this, please set up your environment according to [this page](https://golang.org/doc/code.html#Workspaces).

To build the `gocd` binary, run the `build.sh` command:

```bash
$ ./build.sh
```

This will generate the `gocd` binary in the repository's root directory (i.e., most likely your working directory).

You will likely need recent versions of config-repo plugins to do anything interesting:

```
$ ./gocd configrepo -i yaml.config.plugin fetch # saves to $HOME/.gocd/plugins
$ ./gocd configrepo -i json.config.plugin fetch # saves to $HOME/.gocd/plugins
```

## Usage

There are built-in help screens to the `gocd` binary if you pass in `-h` or invoke with no arguments at all. Will expand later on this.

### Example: Do a syntax check on a config-repo definition file

```
# `-i` specifies the plugin id; in this case, `yaml.config.plugin` is the YAML Config Repo Plugin's identifier.
# Use json.config.plugin for the JSON Config Repo Plugin

$ ./gocd configrepo -i yaml.config.plugin syntax my-pipeline.gocd.yaml
OK
```

### Example: Fetch a config-repo plugin

```
# `-i` specifies the plugin id; in this case, `yaml.config.plugin` is the YAML Config Repo Plugin's identifier.
# Use json.config.plugin for the JSON Config Repo Plugin

$ ./gocd configrepo -i yaml.config.plugin fetch
Downloading https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.8.3/yaml-config-plugin-0.8.3.jar
  Fetched 2.3 MB/2.3 MB (100.0%) complete
```

### Example: Fetch a specific version of a config-repo plugin

```
$ ./gocd configrepo -i yaml.config.plugin fetch --match-version 0.7.0
Downloading https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.7.0/yaml-config-plugin-0.7.0.jar
  Fetched 2.0 MB/2.0 MB (100.0%) complete
```

### Example: Fetch a config-repo plugin matching a version range

```
# the `--match-version` flag accepts a semver (semantic version) query or range string
#
# Examples:
#
# --match-version '>=0.5.0 <0.8.0' # simple range; will resolve to 0.7.0
# --match-version '>=0.5.0 <0.8.0 || >=0.8.1 !0.8.3' # compound range; will resolve to 0.8.2
# --match-version '0.8.x' # wildcard, will match latest 0.8.x release

$ ./gocd configrepo -i yaml.config.plugin fetch --match-version '< 0.7.0'
Downloading https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.6.2/yaml-config-plugin-0.6.2.jar
  Fetched 2.0 MB/2.0 MB (100.0%) complete
```

### Example: Configure settings

```
# Save GoCD server base URL for API calls

$ ./gocd config server-url https://build.gocd.org
```

```
# Save auth credentials (currently, basic auth only) for API calls

$ ./gocd config auth-basic myuser secretpassword
```
