package zoxide

import (
	"runtime"
	"testing"

	zoxidelib "github.com/lazysegtree/go-zoxide"

	"github.com/yorukot/superfile/src/pkg/utils"
)

func setupTestModel() Model {
	return GenerateModel(nil, 50, 80) //nolint:mnd // test dimensions
}

func setupTestModelWithClient(t *testing.T) Model {
	t.Helper()
	zClient, err := zoxidelib.New(zoxidelib.WithDataDir(t.TempDir()))
	if err != nil {
		if runtime.GOOS != utils.OsLinux {
			t.Skipf("Skipping zoxide tests in non-Linux because zoxide client cannot be initialized")
		} else {
			t.Fatalf("zoxide initialization failed")
		}
	}
	return GenerateModel(zClient, 50, 80) //nolint:mnd // test dimensions
}

func setupTestModelWithResults(resultCount int) Model {
	m := setupTestModel()
	m.results = make([]zoxidelib.Result, resultCount)
	for i := range resultCount {
		m.results[i] = zoxidelib.Result{
			Path:  "/test/path" + string(rune('0'+i)),
			Score: float64(100 - i*10), //nolint:mnd // test scores
		}
	}
	return m
}
