package prompt

import (
	"github.com/stretchr/testify/assert"
	"github.com/yorukot/superfile/src/internal/common"
	"testing"
)

func TestModel_getPromptAction(t *testing.T) {
	// Things we want to test
	// 1. No Action
	// 2. Shell command
	// 3. Tokenization failure (Don't test all failures, it will be in tokenize_test.go)
	// 4. Command with incorrect no. of args
	// 5. Invalid command
	// 6. Correct parsing of split, open and cd commands

	testdata := []struct {
		name           string
		text           string
		shellMode      bool
		expectecAction common.ModelAction
		expectedErr    bool
		expectedErrMsg string
	}{
		{
			name:           "No Action",
			text:           "",
			shellMode:      true,
			expectecAction: common.NoAction{},
			expectedErr:    false,
			expectedErrMsg: "",
		},
		{
			name:      "Shell command",
			text:      "abc xyz /def",
			shellMode: true,
			expectecAction: common.ShellCommandAction{
				Command: "abc xyz /def",
			},
			expectedErr:    false,
			expectedErrMsg: "",
		},
		{
			name:           "Tokenization failure",
			text:           "cd ${sdfdsf", // Missing "}"
			shellMode:      false,
			expectecAction: common.NoAction{},
			expectedErr:    true,
			expectedErrMsg: tokenizationError,
		},
		{
			name:           "Split with extra arguments",
			text:           splitCommand + " xyz",
			shellMode:      false,
			expectecAction: common.NoAction{},
			expectedErr:    true,
			expectedErrMsg: splitCommandArgError,
		},
		{
			name:           "cd with 0 arguments",
			text:           cdCommand,
			shellMode:      false,
			expectecAction: common.NoAction{},
			expectedErr:    true,
			expectedErrMsg: "cd command needs exactly one argument, received 0",
		},
		{
			name:           "Invalid command",
			text:           "abcd",
			shellMode:      false,
			expectecAction: common.NoAction{},
			expectedErr:    true,
			expectedErrMsg: "Invalid spf prompt command : abcd",
		},
		{
			name:           "Correct split command",
			text:           splitCommand,
			shellMode:      false,
			expectecAction: common.SplitPanelAction{},
			expectedErr:    false,
			expectedErrMsg: "",
		},
		{
			name:           "Correct cd command",
			text:           cdCommand + " /abc",
			shellMode:      false,
			expectecAction: common.CDCurrentPanelAction{Location: "/abc"},
			expectedErr:    false,
			expectedErrMsg: "",
		},
		{
			name:           "Correct open command",
			text:           openCommand + " /abc",
			shellMode:      false,
			expectecAction: common.OpenPanelAction{Location: "/abc"},
			expectedErr:    false,
			expectedErrMsg: "",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			action, err := getPromptAction(tt.shellMode, tt.text, "/")
			if err != nil {
				assert.True(t, tt.expectedErr)
				cmdErr, ok := err.(invalidCmdError)
				assert.True(t, ok)
				if tt.expectedErrMsg != "" {
					assert.Equal(t, tt.expectedErrMsg, cmdErr.uiMessage())
				}
			}

			assert.Equal(t, tt.expectecAction, action)

		})
	}

}
