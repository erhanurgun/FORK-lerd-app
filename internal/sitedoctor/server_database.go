package sitedoctor

import (
	"fmt"
	"strings"

	"github.com/geodro/lerd/internal/config"
	"github.com/geodro/lerd/internal/serviceops"
)

// listDatabases is the engine-agnostic lookup, hooked so tests can drive the
// check without a running container. The command itself comes from the service
// preset, so no engine SQL is hardcoded here.
var listDatabases = func(service string) ([]string, error) {
	infos, err := serviceops.ListDatabases(service, serviceops.IntrospectCommand(service))
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(infos))
	for _, db := range infos {
		names = append(names, db.Name)
	}
	return names, nil
}

// stubDatabaseLister swaps the lookup for a test and returns the restore func.
func stubDatabaseLister(fn func(string) ([]string, error)) func() {
	prev := listDatabases
	listDatabases = fn
	return func() { listDatabases = prev }
}

// checkServerDatabase fails when the site's database does not exist on the
// engine it points at. Without it such a site 500s on every request while the
// doctor reports nothing wrong: the framework's own migration check cannot reach
// the app either, so it degrades to "couldn't run" and is not counted.
//
// Which database on which service is read from the framework declaration, so a
// project keeping its configuration in a PHP settings file or behind a DSN is
// checked like any other. It used to read DB_CONNECTION, DB_HOST and DB_DATABASE
// by name, which are Laravel's, so the frameworks least likely to be wired that
// way were the ones it could not check.
//
// The fix creates the database rather than migrating it: migrations fail against
// a database the engine does not hold, so the schema has to exist first, and the
// migrate button returns on the re-check that follows.
func checkServerDatabase(path string) (Check, bool) {
	missing, checked := missingDatabases(path)
	if !checked {
		// Either nothing could be asked of an engine, or the project points at no
		// lerd-run database at all: a file database, an external server, or
		// nothing configured. None of those is this check's to judge.
		return Check{}, false
	}
	if len(missing) == 0 {
		return Check{Name: "server_database", Status: StatusOK}, true
	}
	named := make([]string, 0, len(missing))
	for _, t := range missing {
		named = append(named, fmt.Sprintf("%q on %s", t.Database, t.Service))
	}
	return Check{Name: "server_database", Status: StatusFail, Fix: FixCreateDatabase,
		Detail: fmt.Sprintf("%s %s %s not exist. Create %s, then run migrations.",
			plural(len(missing), "Database", "Databases"), strings.Join(named, ", "),
			plural(len(missing), "does", "do"), plural(len(missing), "it", "them"))}, true
}

// MissingDatabases returns the lerd-managed databases a project points at that
// their engine does not hold. Exported so the fix works the set out again
// rather than trusting the client, the same way the service fixes do.
func MissingDatabases(path string) []config.DBTarget {
	missing, _ := missingDatabases(path)
	return missing
}

// missingDatabases pairs the set with whether any engine could be asked at all:
// an unreachable one leaves the site unjudged rather than reported as missing a
// schema that may well exist.
func missingDatabases(path string) (missing []config.DBTarget, checked bool) {
	for _, t := range config.DBTargetsFor(path) {
		names, err := listDatabases(t.Service)
		if err != nil {
			continue
		}
		checked = true
		found := false
		for _, n := range names {
			if strings.EqualFold(n, t.Database) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, t)
		}
	}
	return missing, checked
}
