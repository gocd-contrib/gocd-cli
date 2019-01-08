package cfg

import "fmt"

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

func from(version int) *migration { //
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

func (m *migration) needsMigration(schema dict) bool {
	if _, exists := schema[CONFIG_VERSION]; !exists {
		return false
	}

	if v, ok := schema[CONFIG_VERSION].(int); !ok {
		return false
	} else {
		return v == m.fromVersion
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

	if !m.needsMigration(input) {
		return input, nil
	}

	if output, err := m.body(input); err == nil {
		m.setsVersion(output)
		return output, nil
	} else {
		return nil, err
	}
}
