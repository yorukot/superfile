package processbar

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
)

// TODO: This is duplicated in tests of prompt package, internal package too.
// Fix this code duplication

// Initialize the globals we need for testing
func initGlobals() {
	// Updating globals for test is not a good idea and can lead to all sorts of issue
	// When multiple tests depend on same global variable and want different values
	// Since this is config that would likely stay same, maybe this is okay.
	// Also, this is done in main model's test too.
	// We need to find a better way to do this
	err := common.PopulateGlobalConfigs()
	if err != nil {
		fmt.Printf("error while populating config, err : %v", err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}
	initGlobals()
	m.Run()
}

func TestModelProcessUtils(t *testing.T) {
	m := New()
	p1 := NewProcess("1", "test", OpCopy, 10)
	p2 := NewProcess("2", "test2", OpDelete, 11)

	// ------- Testing AddProcess

	err := m.AddProcess(p1)
	require.NoError(t, err, "Add should succeed without errors")

	err = m.AddProcess(p2)
	require.NoError(t, err, "Add should succeed without errors for second process")

	pRes, ok := m.GetByID(p1.ID)
	require.True(t, ok, "Should be able to get the process we just added")
	assert.Equal(t, p1, pRes, "Should get the correct process value")

	p2Dup := NewProcess("2", "test2_dup", OpCopy, 1)
	err = m.AddProcess(p2Dup)
	var errExp *ProcessAlreadyExistsError
	require.ErrorAs(t, err, &errExp, "Should get ProcessAlreadyExistsError")
	assert.Equal(t, errExp.id, p2Dup.ID, "ID in the error should match with what we sent")

	// ------ Testing AddOrUpdate process
	m.AddOrUpdateProcess(p2Dup)
	pRes, ok = m.GetByID(p2Dup.ID)
	require.True(t, ok)
	assert.Equal(t, p2Dup, pRes, "Should get the correct process value after update")

	p3 := NewProcess("3", "test3", OpExtract, 1)

	// ------ Testing UpdateExisting

	err = m.UpdateExistingProcess(p3)
	var errExpUpdate *NoProcessFoundError
	require.ErrorAs(t, err, &errExpUpdate, "Should get NoProcessFoundError")
	assert.Equal(t, p3.ID, errExpUpdate.id, "ID in the error should match with what we sent")

	assert.True(t, m.HasRunningProcesses())

	// Update all to done
	p1.State = Successful
	p2Dup.Done = p2Dup.Total
	p3.State = Failed
	_ = m.UpdateExistingProcess(p1)
	_ = m.UpdateExistingProcess(p2Dup)
	_ = m.UpdateExistingProcess(p3)

	assert.False(t, m.HasRunningProcesses())
}

func TestModelSetDimenstions(t *testing.T) {
	m := New()

	m.SetDimensions(5, 6)
	assert.Equal(t, 5, m.width, "Correct value should be set")
	assert.Equal(t, 6, m.height, "Correct value should be set")

	m.SetDimensions(minWidth+1, minHeight-1)
	assert.Equal(t, minHeight, m.height, "Min value should be set")
	assert.Equal(t, minWidth+1, m.width, "Given value should be set")
}
