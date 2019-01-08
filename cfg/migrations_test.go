package cfg

import (
	"errors"
	"fmt"
	"testing"
)

func TestMigrations(t *testing.T) {
	as := asserts(t)
	c, err := makeConf(fmt.Sprintf(`%s: %d
foo: bar
`, CONFIG_VERSION, CURRENT_VERSION-2))
	as.ok(err)

	migs := []*migration{
		from(CURRENT_VERSION-2).
			to(CURRENT_VERSION-1).
			do(`change foo to baz`, func(cfg dict) (dict, error) {
				if v, ok := cfg[`foo`]; ok {
					cfg[`baz`] = v
					delete(cfg, `foo`)
				}
				return cfg, nil
			}),

		from(CURRENT_VERSION-1).
			to(CURRENT_VERSION).
			do(`baz becomes a list`, func(cfg dict) (dict, error) {
				if v, ok := cfg[`baz`]; ok {
					cfg[`baz`] = []interface{}{v}
				}
				return cfg, nil
			}),
	}

	as.ok(c.Migrate(migs))
	as.configEq(dict{
		`config_version`: CURRENT_VERSION,
		`baz`:            []string{`bar`},
	}, c.fs)
}

func TestMigrationsDoesNotModifyConfigFileOnError(t *testing.T) {
	as := asserts(t)
	c, err := makeConf(fmt.Sprintf(`%s: %d
foo: bar
`, CONFIG_VERSION, CURRENT_VERSION-2))
	as.ok(err)

	migs := []*migration{
		from(CURRENT_VERSION-2).
			to(CURRENT_VERSION-1).
			do(`change foo to baz`, func(cfg dict) (dict, error) {
				if v, ok := cfg[`foo`]; ok {
					cfg[`baz`] = v
					delete(cfg, `foo`)
				}
				return cfg, nil
			}),

		from(CURRENT_VERSION-1).
			to(CURRENT_VERSION).
			do(`should blow up`, func(cfg dict) (dict, error) {
				return nil, errors.New(`boom!`)
			}),
	}

	as.err(`boom!`, c.Migrate(migs))

	// should not change what is on disk
	as.configEq(dict{
		CONFIG_VERSION: CURRENT_VERSION - 2,
		`foo`:          `bar`,
	}, c.fs)
}

func TestOnlyMigratesWhenPrerequisiteVersionMet(t *testing.T) {
	as := asserts(t)
	c, err := makeConf(fmt.Sprintf(`%s: %d
foo: bar
`, CONFIG_VERSION, CURRENT_VERSION-1))
	as.ok(err)

	didRun := false
	migs := []*migration{
		from(CURRENT_VERSION-2).
			to(CURRENT_VERSION-1).
			do(`should not run`, func(cfg dict) (dict, error) {
				return nil, errors.New(`should not have run`)
			}),

		from(CURRENT_VERSION-1).
			to(CURRENT_VERSION).
			do(`should only run this`, func(cfg dict) (dict, error) {
				didRun = true
				return cfg, nil
			}),
	}

	as.ok(c.Migrate(migs))
	as.is(didRun)
}

func TestReturnsErrorWhenVersionTooNew(t *testing.T) {
	as := asserts(t)
	c, err := makeConf(fmt.Sprintf(`%s: %d
foo: bar
`, CONFIG_VERSION, CURRENT_VERSION+1))
	as.ok(err)

	didRun := false
	migs := []*migration{
		from(CURRENT_VERSION).to(CURRENT_VERSION+1).do(`should not run`, func(dict) (dict, error) {
			didRun = true
			return nil, nil
		}),
	}

	as.not(didRun)
	as.err(
		fmt.Sprintf(`%q: %d is not supported in this CLI version`, CONFIG_VERSION, CURRENT_VERSION+1),
		c.Migrate(migs),
	)
}
