package configrepo

import (
	"fmt"
	"strings"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/materials"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GitRepoCmd(doCreate bool) (cmd *cobra.Command) {
	var (
		repoUrl string
		branch  string

		builder = &RepoRunner{m: materials.NewGit(), props: make(PropertySet)}
	)

	if doCreate {
		cmd = &cobra.Command{
			Use:   `git <id> --url=<url>`,
			Short: `Creates a git config-repo`,
			Args:  cobra.ExactArgs(1),
			Run: builder.AddRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				utils.AbortLoudlyOnError(m.SetRequiredString(`url`, repoUrl))
				m.SetStringWithDefault(`branch`, branch, `master`)
			}),
		}
	} else {
		cmd = &cobra.Command{
			Use:   `git <id>`,
			Short: `Updates a git config-repo`,
			Args:  cobra.ExactArgs(1),
			Run: builder.UpdateRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				m.SetStringIfFlagSet(`url`, repoUrl, cmd.Flag(`url`))
				m.SetStringIfFlagSet(`branch`, branch, cmd.Flag(`branch`))
			}),
		}
	}

	return builder.SetupFlags(cmd, func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&repoUrl, `url`, ``, `Set the git repository url`)
		cmd.Flags().StringVar(&branch, `branch`, ``, `Set the git branch name`)

		if doCreate {
			cmd.MarkFlagRequired(`url`)
		}
	})
}

func HgRepoCmd(doCreate bool) (cmd *cobra.Command) {
	var (
		repoUrl string

		builder = &RepoRunner{m: materials.NewHg(), props: make(PropertySet)}
	)

	if doCreate {
		cmd = &cobra.Command{
			Use:   `hg <id> --url=<url>`,
			Short: `Creates an hg config-repo`,
			Args:  cobra.ExactArgs(1),
			Run: builder.AddRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				utils.AbortLoudlyOnError(m.SetRequiredString(`url`, repoUrl))
			}),
		}
	} else {
		cmd = &cobra.Command{
			Use:   `hg <id>`,
			Short: `Updates an hg config-repo`,
			Args:  cobra.ExactArgs(1),
			Run: builder.UpdateRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				m.SetStringIfFlagSet(`url`, repoUrl, cmd.Flag(`url`))
			}),
		}
	}

	return builder.SetupFlags(cmd, func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&repoUrl, `url`, ``, `Set the hg repository url`)

		if doCreate {
			cmd.MarkFlagRequired(`url`)
		}
	})
}

func SvnRepoCmd(doCreate bool) (cmd *cobra.Command) {
	var (
		repoUrl string
		extern  = true

		builder = &RepoRunner{m: materials.NewSvn(), props: make(PropertySet)}
	)

	if doCreate {
		cmd = &cobra.Command{
			Use:   `svn <id> --url=<url>`,
			Args:  cobra.ExactArgs(1),
			Short: `Creates an svn config-repo`,
			Run: builder.AddRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				utils.AbortLoudlyOnError(m.SetRequiredString(`url`, repoUrl))
				m.SetBool(`check_externals`, extern)
			}),
		}
	} else {
		cmd = &cobra.Command{
			Use:   `svn <id>`,
			Args:  cobra.ExactArgs(1),
			Short: `Updates an svn config-repo`,
			Run: builder.UpdateRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				utils.AbortLoudlyOnError(m.SetRequiredString(`url`, repoUrl))
				m.SetStringIfFlagSet(`url`, repoUrl, cmd.Flag(`url`))
				m.SetBoolIfFlagSet(`check_externals`, extern, cmd.Flag(`check-externals`))
			}),
		}
	}

	return builder.SetupFlags(cmd, func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&repoUrl, `url`, ``, `Set the svn repository url`)
		cmd.Flags().BoolVar(&extern, `check-externals`, true, `Config-repo should check svn externals (set =false to disable)`)

		if doCreate {
			cmd.MarkFlagRequired(`url`)
		}
	})
}

func P4RepoCmd(doCreate bool) (cmd *cobra.Command) {
	var (
		hostPort string
		view     string
		useTix   = true

		builder = &RepoRunner{m: materials.NewP4(), props: make(PropertySet)}
	)

	if doCreate {
		cmd = &cobra.Command{
			Use:   `p4 <id> --host-port=<[scheme:]host:port> --view=<view>`,
			Args:  cobra.ExactArgs(1),
			Short: `Creates a p4 config-repo`,
			Run: builder.AddRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				utils.AbortLoudlyOnError(m.SetRequiredString(`port`, hostPort))
				utils.AbortLoudlyOnError(m.SetRequiredString(`view`, view))
				m.SetBool(`use_tickets`, useTix)
			}),
		}
	} else {
		cmd = &cobra.Command{
			Use:   `p4 <id>`,
			Args:  cobra.ExactArgs(1),
			Short: `Updates a p4 config-repo`,
			Run: builder.UpdateRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				m.SetStringIfFlagSet(`port`, hostPort, cmd.Flag(`host-port`))
				m.SetStringIfFlagSet(`view`, view, cmd.Flag(`view`))
				m.SetBoolIfFlagSet(`use_tickets`, useTix, cmd.Flag(`use-tickets`))
			}),
		}
	}

	return builder.SetupFlags(cmd, func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&hostPort, `host-port`, ``, "Set the p4 repository `[scheme:]host:port`")
		cmd.Flags().StringVar(&view, `view`, ``, `Set the p4 repository view`)
		cmd.Flags().BoolVar(&useTix, `use-tickets`, true, "Config-repo should use p4 ticket-based authentication; set =false to\n  disable")

		if doCreate {
			cmd.MarkFlagRequired(`host-port`)
			cmd.MarkFlagRequired(`view`)
		}
	})
}

func TfsRepoCmd(doCreate bool) (cmd *cobra.Command) {
	var (
		repoUrl string
		proj    string
		domain  string

		builder = &RepoRunner{m: materials.NewTfs(), props: make(PropertySet)}
	)

	if doCreate {
		cmd = &cobra.Command{
			Use:   `tfs <id> --url=<url> --project-path=<project-path>`,
			Args:  cobra.ExactArgs(1),
			Short: `Creates a tfs config-repo`,
			Run: builder.AddRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				utils.AbortLoudlyOnError(m.SetRequiredString(`url`, repoUrl))
				utils.AbortLoudlyOnError(m.SetRequiredString(`project_path`, proj))
				m.SetStringIfFlagSet(`domain`, domain, cmd.Flag(`domain`))
			}),
		}
	} else {
		cmd = &cobra.Command{
			Use:   `tfs <id>`,
			Args:  cobra.ExactArgs(1),
			Short: `Updates a tfs config-repo`,
			Run: builder.UpdateRepo(func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material) {
				m.SetStringIfFlagSet(`url`, repoUrl, cmd.Flag(`url`))
				m.SetStringIfFlagSet(`project_path`, proj, cmd.Flag(`project-path`))
				m.SetStringIfFlagSet(`domain`, domain, cmd.Flag(`domain`))
			}),
		}
	}

	return builder.SetupFlags(cmd, func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&repoUrl, `url`, ``, "Set the tfs repository url")
		cmd.Flags().StringVar(&proj, `project-path`, ``, `Set the tfs repository project-path`)
		cmd.Flags().StringVar(&domain, `domain`, ``, `Set the tfs repository domain`)

		if doCreate {
			cmd.MarkFlagRequired(`url`)
			cmd.MarkFlagRequired(`project-path`)
		}
	})
}

type RepoRunner struct {
	thisRepo            *ConfigRepo
	m                   materials.Material
	props               PropertySet
	user, pass, encPass string
}

// Configures common flags across repos, including user/pass and property setting flags
// Takes an init() function to set up the remainder of the flags
func (r *RepoRunner) SetupFlags(cmd *cobra.Command, init func(cmd *cobra.Command)) *cobra.Command {
	init(cmd)

	if r.supportsAuthFlags() {
		r.addAuthFlags(cmd)
	}
	r.addPropertyFlags(cmd)

	cmd.Flags().SortFlags = false // print them in the order defined
	return cmd
}

func (r *RepoRunner) AddRepo(apply func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		DieWhenPluginIdNotSet()

		id := args[0]

		apply(cmd, r.thisRepo, r.m)

		if r.supportsAuthFlags() {
			r.m.SetStringIfFlagSet(`username`, r.user, cmd.Flag(`user`))
			if !isFlagSet(cmd, `encrypted-password`) {
				r.m.SetStringIfFlagSet(`password`, r.pass, cmd.Flag(`password`))
			}
			r.m.SetStringIfFlagSet(`encrypted_password`, r.encPass, cmd.Flag(`encrypted-password`))
		}

		r.thisRepo = NewConfigRepo(id, PluginId, r.m, r.props.Value()...)

		utils.AbortLoudlyOnError(Model.AddRepo(r.thisRepo, func(repo *ConfigRepo) error {
			utils.Echofln(`Created repo %q`, repo.Id)
			return utils.InspectError(api.PrettyPrintJson(repo), `printing config-repo on successful config-repo create`)
		}))
	}
}

func (r *RepoRunner) UpdateRepo(apply func(cmd *cobra.Command, repo *ConfigRepo, m materials.Material)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		DieWhenPluginIdNotSet()
		DieWhenNoFlagsSetForUpdate(cmd)

		id := args[0]

		utils.Debug(`Need to fetch ETag for config-repo %q before update`, id)
		utils.AbortLoudlyOnError(Model.FetchRepo(id, func(repo *ConfigRepo) error {
			if repo.Material.Type() != r.m.Type() {
				return fmt.Errorf(`Expected config-repo %q to use a %s material, but was %s`, id, r.m.Type(), repo.Material.Type())
			}

			r.thisRepo = repo
			r.m = repo.Material
			return nil
		}))

		apply(cmd, r.thisRepo, r.m)

		if r.supportsAuthFlags() {
			r.m.SetStringIfFlagSet(`username`, r.user, cmd.Flag(`user`))
			if !isFlagSet(cmd, `encrypted-password`) {
				r.m.SetStringIfFlagSet(`password`, r.pass, cmd.Flag(`password`))
			}
			r.m.SetStringIfFlagSet(`encrypted_password`, r.encPass, cmd.Flag(`encrypted-password`))
		}

		if isFlagSet(cmd, `property`) || isFlagSet(cmd, `secret-property`) {
			r.thisRepo.Configuration = r.props.Value()
		}

		utils.AbortLoudlyOnError(Model.UpdateRepo(r.thisRepo, func(repo *ConfigRepo) error {
			return utils.InspectError(api.PrettyPrintJson(repo), `printing config-repo on successful config-repo create`)
		}))
	}
}

func (r *RepoRunner) supportsAuthFlags() bool {
	switch r.m.Type() {
	case `git`:
		return false
	case `hg`:
		return false
	default:
		return true
	}
}

func (r *RepoRunner) addAuthFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&r.user, `user`, ``, fmt.Sprintf(`Set the %s repository user`, r.m.Type()))
	cmd.Flags().StringVar(&r.pass, `password`, ``, fmt.Sprintf(`Set the %s repository password`, r.m.Type()))
	cmd.Flags().StringVar(&r.encPass, `encrypted-password`, ``, strings.Join([]string{
		`Same as --password, but assumes the value is the encrypted form;`,
		`  always takes precedence over --password`,
	}, "\n"))
}

func (r *RepoRunner) addPropertyFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(MakePropertySetValue(&r.props, NewPlainTextProperty), `property`, `w`, strings.Join([]string{
		`Set plugin-specific configuration options; may be repeated additively,`,
		`  but for duplicate keys, the last value wins`,
	}, "\n"))
	cmd.Flags().VarP(MakePropertySetValue(&r.props, NewSecretProperty), `secret-property`, `W`, strings.Join([]string{
		"Same as --property, but treats the value as an \"encrypted_value\"; this",
		`  option DOES NOT perform encryption -- you must do that on your own`,
	}, "\n"))
}

func isFlagSet(cmd *cobra.Command, key string) bool {
	return cmd.Flag(key).Changed
}

func anyFlagsSet(cmd *cobra.Command) bool {
	set := false
	cmd.Flags().Visit(func(f *pflag.Flag) {
		set = set || f.Changed
	})
	return set
}

func DieWhenNoFlagsSetForUpdate(cmd *cobra.Command) {
	if !anyFlagsSet(cmd) {
		utils.DieLoudly(1, `You haven't specified any settings to update`)
	}
}
