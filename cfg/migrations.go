package cfg

import (
	"fmt"

	"github.com/gocd-contrib/gocd-cli/utils"
)

// as the config format evolves, we will need to add config migrations here.
// each migration is responsible for checking the prerequisite version (and
// exiting early if not applicable) and bumping the version for the next
// migration.
//
// Constructing a migration:
//
//   from(1).to(2).do(`Name of migration`, func(schema dict) (dict, error) {
//     ... perform schema changes here ...
//   })
var migrations = []*migration{}

func applyMigrations(initial dict, migrations []*migration) (migrated dict, err error) {
	migrated = initial

	for _, mig := range migrations {
		if migrated, err = mig.run(migrated); err != nil {
			migrated = nil
			return
		}
	}
	return
}

type migration struct {
	fromVersion int
	toVersion   int
	name        string
	body        func(dict) (dict, error)
}

func from(version int) *migration {
	return &migration{fromVersion: version}
}

func (m *migration) to(version int) *migration {
	m.toVersion = version
	return m
}

func (m *migration) do(name string, body func(dict) (dict, error)) *migration {
	m.name = name
	m.body = body
	return m
}

func (m *migration) needsMigration(schema dict) (bool, error) {
	if _, exists := schema[CONFIG_VERSION]; !exists {
		return false, fmt.Errorf(`Config file is missing %q key`, CONFIG_VERSION)
	}

	if v, ok := schema[CONFIG_VERSION].(int); !ok {
		return false, fmt.Errorf(`Value for key %q must be numeric; instead, was of type: %T`, CONFIG_VERSION, v)
	} else {
		return v == m.fromVersion, nil
	}
}

func (m *migration) setsVersion(schema dict) {
	schema[CONFIG_VERSION] = m.toVersion
}

func (m *migration) run(input dict) (dict, error) {
	if `` == m.name {
		return nil, fmt.Errorf(`Migration %d -> %d has no name`, m.fromVersion, m.toVersion)
	}

	if nil == m.body {
		return nil, fmt.Errorf(`Migration %q has no body`, m.name)
	}

	utils.Debug(`Need to apply migration? => %q`, m.name)
	if outdated, err := m.needsMigration(input); err != nil {
		return nil, err
	} else {
		if !outdated { // up-to-date, no migration needed
			utils.Debug(`Config is already at version %d or better; skipping...`, m.toVersion)
			return input, nil
		}
	}

	utils.Debug(`Applying %q...`, m.name)
	if output, err := m.body(input); err == nil {
		m.setsVersion(output)
		utils.Debug(`Successfully applied %q, now at version %d`, m.name, m.toVersion)
		return output, nil
	} else {
		utils.Debug(`Failed to apply %q`, m.name)
		return nil, err
	}
}
