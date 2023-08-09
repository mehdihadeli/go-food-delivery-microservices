package testfixture

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/go-testfixtures/testfixtures/v3"
)

func RunPostgresFixture(db *sql.DB, fixturePaths []string, data map[string]interface{}) error {
	// determine the project's root path
	_, callerPath, _, _ := runtime.Caller(1) // nolint:dogsled

	// Root folder of this project
	rootPath := filepath.Join(filepath.Dir(callerPath), "../../../..")

	// assemble a list of fixtures paths to be loaded
	for i := range fixturePaths {
		fixturePaths[i] = fmt.Sprintf("%v/%v", rootPath, filepath.ToSlash(fixturePaths[i]))
	}

	// https://github.com/go-testfixtures/testfixtures
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Template(),
		testfixtures.TemplateData(data),
		// Paths must come after Template() and TemplateData()
		testfixtures.Paths(fixturePaths...),
	)
	if err != nil {
		return err
	}

	// load fixtures into DB
	return fixtures.Load()
}
