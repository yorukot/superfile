# prompt package
This is for the Prompt modal of superfile

Handles user input updates, spf model updates, and returns a PromptAction to model. 


# Coverage

```bash
cd /path/to/ui/prompt
# Basic coverage
go test -cover

# HTML report
go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```
Current coverage is 74.0%.

# To-dos
- Use bubbltea's textinput Model's .prompt field for '>' or ':' prompt