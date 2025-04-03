package prompt

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/yorukot/superfile/src/internal/common/utils"
	"testing"
)

func Test_tokenizePromptCommand(t *testing.T) {
	// Just test that we can split as expected
	// Empty string, trailing and leading whitespace
	// single, and multiple tokens
	// special characters
	// Dont try to test shell substitution in this. This is just
	// to test that tokenize function can handle the results of shell
	// substitution as expected

	defaultCwd := "/"
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
			res, err := tokenizePromptCommand(tt.command, defaultCwd)
			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.isErrorExpected, err != nil)
		})
	}
}

func Test_resolveShellSubstitution(t *testing.T) {
	// We want to test
	// Empty string
	// Strings without $ - Normal, trailing and leading whitespace, special char
	// also with $ but no ${} or $() - $HOME , _$$_xyz
	// Ill formatted substitutions - missing '$', missing ')' or '}' or '{'
	// Multiple correct ${} and $()
	// Empty ${} and $()
	// ${ $() } -> Should not work , $(echo ${} $(echo $())) -> Should work
	// cd $(echo $(echo hi))
	// no output shell commands $(true), newline output $(echo -e "\n")
	// env var - not found
	// substitution command times out
	defaultCwd := "/"
	utils.SetRootLoggerToStdout(true)
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
			errorToMatch:    bracketParMatchError(),
		},
		{
			name:            "Ill formed command 2",
			command:         "abc $(abc) syt ${ sdfc ( {)}",
			expectedResult:  "",
			isErrorExpected: true,
			errorToMatch:    curlyBracketParMatchError(),
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
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveShellSubstitution(shellSubTimeoutInTests, tt.command, defaultCwd)
			assert.Equal(t, tt.expectedResult, result)
			if err != nil {
				assert.True(t, tt.isErrorExpected)
			}
		})
	}

	t.Run("Testing shell substitution timeout", func(t *testing.T) {
		result, err := resolveShellSubstitution(shellSubTimeoutInTests, "$(sleep 0.1)", defaultCwd)
		assert.Empty(t, result)
		assert.NotNil(t, err)
		assert.ErrorAs(t, err, &context.DeadlineExceeded)
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
			res := findEndingParenthesis([]rune(tt.value), tt.openIdx, tt.openPar, tt.closePar)
			assert.Equal(t, tt.expectedRes, res)
		})
	}
}
