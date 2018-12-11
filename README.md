# gocd-cli

A command-line companion to GoCD. In its current state, it can only perform syntax checking for configrepo definitions. More to come later.

## Development

This is a [golang](https://golang.org/) project, so if you are new to this, please set up your environment according to [this page](https://golang.org/doc/code.html#Workspaces).

To build the `gocd` binary, run the `build.sh` command:

```bash
$ ./build.sh
```

This will generate the `gocd` binary in the root folder.

```bash
# create the plugins folder
mkdir -p $HOME/.gocd/plugins
```

Download the latest release of the [JSON](https://github.com/tomzo/gocd-json-config-plugin/releases) and/or [YAML](https://github.com/tomzo/gocd-yaml-config-plugin/releases) plugins and save them to `$HOME/.gocd/plugins`

## Usage

There are built-in help screens to the `gocd` binary if you pass in `-h` or invoke with no arguments at all. Will expand later on this.

### Example to do a syntax check on a configrepo definition file

```
# `-i` specifies the plugin id; in this case, `yaml.config.plugin` is the YAML Config Repo Plugin's identifier.
# Use json.config.plugin for the JSON Config Repo Plugin

$ ./gocd configrepo check -i yaml.config.plugin my-pipeline.gocd.yaml
$ OK
```
