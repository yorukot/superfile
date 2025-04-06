package prompt

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	spfTestEnvVar1 = "SPF_TEST_ENV_VAR1"
	spfTestEnvVar2 = "SPF_TEST_ENV_VAR2"
	spfTestEnvVar3 = "SPF_TEST_ENV_VAR3"
	spfTestEnvVar4 = "SPF_TEST_ENV_VAR4"
)

var testEnvValues = map[string]string{
	spfTestEnvVar1: "1",
	spfTestEnvVar2: "hello",
	spfTestEnvVar3: "",
}

func Test_tokenizePromptCommand(t *testing.T) {
	// Just test that we can split as expected
	// Don't try to test shell substitution in this. This is just
	// to test that tokenize function can handle the results of shell
	// substitution as expected

	testdata := []struct {
		name            string
		command         string
		expectedRes     []string
		isErrorExpected bool
	}{
		{
			name:            "Empty String",
			command:         "",
			expectedRes:     []string{},
			isErrorExpected: false,
		},
		{
			name:            "Parenthesis issue",
			command:         "abcd $(xyz",
			expectedRes:     nil,
			isErrorExpected: true,
		},
		{
			name:            "Parenthesis issue - But no dollar",
			command:         "abcd (xyz",
			expectedRes:     []string{"abcd", "(xyz"},
			isErrorExpected: false,
		},
		{
			name:            "Whitespace",
			command:         "    a b  c  ",
			expectedRes:     []string{"a", "b", "c"},
			isErrorExpected: false,
		},
		{
			name:            "Single token",
			command:         "()",
			expectedRes:     []string{"()"},
			isErrorExpected: false,
		},
		{
			name:            "Special characters",
			command:         "() \t\n\t a $5^&*\v\a\n\uF0AC",
			expectedRes:     []string{"()", "a", "$5^&*", "\a", "\uF0AC"},
			isErrorExpected: false,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tokenizePromptCommand(tt.command, defaultTestCwd)
			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.isErrorExpected, err != nil)
		})
	}
}

func Test_resolveShellSubstitution(t *testing.T) {
	testdata := []struct {
		name            string
		command         string
		expectedResult  string
		isErrorExpected bool
		errorToMatch    error
	}{
		// Test with no substitution being performed
		{
			name:            "Empty String",
			command:         "",
			expectedResult:  "",
			isErrorExpected: false,
			errorToMatch:    nil,
		},
		{
			name:            "String without substitution requirement",
			command:         "   a b c $%^ () {} \a\v\t \u0087",
			expectedResult:  "   a b c $%^ () {} \a\v\t \u0087",
			isErrorExpected: false,
			errorToMatch:    nil,
		},
		{
			name:            "Ill formed command 1",
			command:         "abc $(abc",
			expectedResult:  "",
			isErrorExpected: true,
			errorToMatch:    roundBracketMatchError(),
		},
		{
			name:            "Ill formed command 2",
			command:         "abc $(abc) syt ${ sdfc ( {)}",
			expectedResult:  "",
			isErrorExpected: true,
			errorToMatch:    curlyBracketMatchError(),
		},

		// Test with substitution being performed
		{
			name:            "Basic substitution",
			command:         "$(echo abc)",
			expectedResult:  "abc\n",
			isErrorExpected: false,
			errorToMatch:    nil,
		},
		// Might not work on windows ?
		{
			name:            "Command with internal substitution",
			command:         "$(echo $(echo abc))",
			expectedResult:  "abc\n",
			isErrorExpected: false,
			errorToMatch:    nil,
		},
		{
			name:            "Multiple substitution",
			command:         fmt.Sprintf("$(echo $(echo $%s)) ${%s}", spfTestEnvVar1, spfTestEnvVar2),
			expectedResult:  fmt.Sprintf("%s\n %s", testEnvValues[spfTestEnvVar1], testEnvValues[spfTestEnvVar2]),
			isErrorExpected: false,
			errorToMatch:    nil,
		},
		{
			name:            "Non Existing env var",
			command:         fmt.Sprintf("${%s}", spfTestEnvVar4),
			expectedResult:  "",
			isErrorExpected: true,
			errorToMatch:    envVarNotFoundError{varName: spfTestEnvVar4},
		},
		{
			name:            "Shell substitution inside env var substitution",
			command:         "${$(pwd)}",
			expectedResult:  "",
			isErrorExpected: true,
			errorToMatch:    envVarNotFoundError{varName: "$(pwd)"},
		},
		{
			name:            "Empty output",
			command:         "cd abc $(true)",
			expectedResult:  "cd abc ",
			isErrorExpected: false,
			errorToMatch:    nil,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveShellSubstitution(shellSubTimeoutInTests, tt.command, defaultTestCwd)
			assert.Equal(t, tt.expectedResult, result)
			if err != nil {
				assert.True(t, tt.isErrorExpected)
				if tt.errorToMatch != nil {
					assert.ErrorIs(t, err, tt.errorToMatch)
				}
			}
		})
	}

	t.Run("Testing shell substitution timeout", func(t *testing.T) {
		result, err := resolveShellSubstitution(shellSubTimeoutInTests, "$(sleep 0.1)", defaultTestCwd)
		assert.Empty(t, result)
		require.Error(t, err)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}

func Test_findEndingParenthesis(t *testing.T) {

	testdata := []struct {
		name        string
		value       string
		openIdx     int
		openPar     rune
		closePar    rune
		expectedRes int
	}{
		{
			name:        "Empty String",
			value:       "",
			openIdx:     0,
			openPar:     '(',
			closePar:    ')',
			expectedRes: -1,
		},
		{
			name:        "Invalid input",
			value:       "abc",
			openIdx:     0,
			openPar:     '(',
			closePar:    ')',
			expectedRes: -1,
		},
		{
			name:        "Simple",
			value:       "abc(def)",
			openIdx:     3,
			openPar:     '(',
			closePar:    ')',
			expectedRes: 7,
		},
		{
			name:  "Nesting Example 1",
			value: "abc(d(e{f})gh)",
			//------01234567890123
			openIdx:     3,
			openPar:     '(',
			closePar:    ')',
			expectedRes: 13,
		},
		{
			name:  "Nesting Example 2",
			value: "abc(d(e{f})gh)",
			//------01234567890123
			openIdx:     5,
			openPar:     '(',
			closePar:    ')',
			expectedRes: 10,
		},
		{
			name:  "Nesting Example 2",
			value: "abc(d(e{f(x}))gh)",
			//------01234567890123456
			openIdx:     7,
			openPar:     '{',
			closePar:    '}',
			expectedRes: 11,
		},
		{
			name:  "No Closing Parenthesis 1",
			value: "abc(def}",
			//------012345678901234
			openIdx:     3,
			openPar:     '(',
			closePar:    ')',
			expectedRes: 8,
		},
		{
			name:  "No Closing Parenthesis 2",
			value: "abc((d(e{f})gh)",
			//------012345678901234
			openIdx:     3,
			openPar:     '(',
			closePar:    ')',
			expectedRes: 15,
		},
		{
			name:  "Asymmetric Parenthesis",
			value: "abc((d(e{f}>gh)",
			//------012345678901234
			openIdx:     8,
			openPar:     '{',
			closePar:    '>',
			expectedRes: 11,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			res := findEndingBracket([]rune(tt.value), tt.openIdx, tt.openPar, tt.closePar)
			assert.Equal(t, tt.expectedRes, res)
		})
	}
}
