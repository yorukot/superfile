# prompt package
This is for the Prompt modal of superfile

Handles user input updates, spf model updates, and returns a PromptAction to model. 

# Todos
x Hints rendering via prefix match
x Handling updates from spf model
x Allow both : and > keys
x Finish todos of prompt
x Take config for exit_after_success
x Remove hardcoded constants
x Run linter on this with hard settings
x Printing return code from shell command
x Check github PR todos
x Make $() and ${} work !!!!
x Do we need textInput.showsuggestions ? No
~ Fix Rendering code
~ Superfile file prompt should appear in the middle of the whole terminal (vertically)
~ Prompt gets resized to smaller value based on total width. and should not break.
~ Make UI looks pretty
x Unit tests
x Check html of unit test coverage

x No - Benchmark test
x No - Fuzzing test
x Testmain
x -debug arg

x Ask on discord for suggestions


# PR Todos
x Remove all new todos
x No More refactoring ? of stuff to common/utils ?
x Bigger unit testcase for model's prompt feature.
x Review go module design docs/videos and see if current design is okay
~ Testsuite test - Later
x Remove hardcoded constants
x Unit tests for other utils that you added

- Coderabbit review and fix comments
- Ask coderabbit for typos, or hardcoding, etc.
- Test on windows
- Execute unit test on windows
- Self Code Review
- More self sanity testing
- Github issue for next AIs
- Merge to a branch

# Coverage
```bash
cd /path/to/ui/prompt
# Basic coverage
go test -cover

# HTML report
go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```