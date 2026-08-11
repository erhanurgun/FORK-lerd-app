package cli

import (
	"testing"

	"github.com/geodro/lerd/internal/config"
)

// Which databases a framework can use is the definition's to declare, like
// everything else about how it is wired. Offering SQLite to a framework that
// cannot use it is how a project ends up picking a database its application
// will never open.
func TestBuildDatabaseOptions_offersSQLiteOnlyWhereDeclared(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_DATA_HOME", t.TempDir())

	declares := &config.Framework{Env: config.FrameworkEnvConf{
		SQLite: &config.FrameworkServiceDef{Vars: []string{"DB_CONNECTION=sqlite"}},
		Services: map[string]config.FrameworkServiceDef{
			"mysql": {Vars: []string{"DB_HOST=lerd-mysql"}},
		},
	}}
	_, names := buildDatabaseOptions(declares)
	if !names["sqlite"] {
		t.Error("a framework declaring sqlite was not offered it")
	}

	declaresNot := &config.Framework{Env: config.FrameworkEnvConf{Services: map[string]config.FrameworkServiceDef{
		"mysql": {Vars: []string{"DB_HOST=lerd-mysql"}},
	}}}
	if _, names := buildDatabaseOptions(declaresNot); names["sqlite"] {
		t.Error("a framework that declares no sqlite wiring was offered it anyway")
	}
}

// A file database is not a service: nothing installs it, starts it or draws a
// card for it. A definition that names it among the services is not what
// declares it, and an older binary reading that entry would try to start it.
func TestBuildDatabaseOptions_ignoresSQLiteAmongServices(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_DATA_HOME", t.TempDir())

	asService := &config.Framework{Env: config.FrameworkEnvConf{Services: map[string]config.FrameworkServiceDef{
		"sqlite": {Vars: []string{"DB_CONNECTION=sqlite"}},
		"mysql":  {Vars: []string{"DB_HOST=lerd-mysql"}},
	}}}
	if _, names := buildDatabaseOptions(asService); names["sqlite"] {
		t.Error("sqlite listed as a service was taken as a declaration")
	}
}

// Which keys point a project at its file database is the definition's to say.
// Writing DB_CONNECTION into a CakePHP project reaches nothing it reads.
func TestSQLiteVarsFor_comeFromTheDefinition(t *testing.T) {
	cakeish := &config.Framework{Env: config.FrameworkEnvConf{SQLite: &config.FrameworkServiceDef{
		Vars: []string{`Datasources.default.driver=Cake\Database\Driver\Sqlite`},
	}}}
	got := sqliteVarsFor(cakeish)
	if len(got) != 1 || got[0] != `Datasources.default.driver=Cake\Database\Driver\Sqlite` {
		t.Errorf("declared sqlite vars were not used, got %v", got)
	}

	if got := sqliteVarsFor(&config.Framework{}); len(got) == 0 || got[0] != "DB_CONNECTION=sqlite" {
		t.Errorf("a framework declaring none did not fall back to the dotenv keys, got %v", got)
	}
}

// A project lerd recognises no framework for keeps the option: nothing has
// declared otherwise, and a file database is a reasonable answer for it.
func TestBuildDatabaseOptions_keepsSQLiteWithoutAFramework(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_DATA_HOME", t.TempDir())

	if _, names := buildDatabaseOptions(nil); !names["sqlite"] {
		t.Error("a project with no framework was not offered sqlite")
	}
	empty := &config.Framework{}
	if _, names := buildDatabaseOptions(empty); !names["sqlite"] {
		t.Error("a framework declaring no env services at all was not offered sqlite")
	}
}
