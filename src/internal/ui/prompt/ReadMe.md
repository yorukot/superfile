# prompt package
This is for the Prompt modal of superfile

Handles user input updates, spf model updates, and returns a PromptAction to model. 

# Todos
x Hints rendering via prefix match
x Handling updates from spf model
x Allow both : and > keys
- Finish todos of prompt
x Take config for exit_after_success
x Remove hardcoded constants
x Run linter on this with hard settings
x Printing return code from shell command
x Check github PR todos
x Make $() and ${} work !!!!
- Do we need textInput.showsuggestions ?
- Superfile file prompt should appear in the middle of the whole terminal (vertically)
- Prompt gets resized to smaller value based on total width. and should not break.
- Implement ... wrapping ?
- Unit tests
- Check html of unit test coverage

- Benchmark test
- Fuzzing test
- Testmain
- -debug arg

- Make UI looks pretty
- Ask on discord for suggestions


# PR Todos
- Remove all new todos
x No More refactoring ? of stuff to common/utils ?
x Bigger unit testcase for model's prompt feature.
- Review go module design docs/videos and see if current design is okay
- Testsuite test
x Remove hardcoded constants
- Unit tests for other utils that you added

- Coderabbit review and fix comments
- Ask coderabbit for typos, or hardcoding, etc.
- Test on windows
- Execute unit test and testsuite on windows
- Self Code Review
- More self sanity testing

# Coverage
```bash
cd /path/to/ui/prompt
# Basic coverage
go test -cover

# HTML report
go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```