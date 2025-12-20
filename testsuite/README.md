## Coding style rules
- Prefer using strong typing
- Prefer using type hinting for the first time the variable is declared, and for functions paremeters and return types
- Use `-> None` to explicitly indicate no return value

### Ideas
- Recommended to integrate your IDE with PEP8 to highlight PEP8 violations in real-time
- Enforcing PEP8 via `pylint flake8 pycodestyle` and via pre commit hooks

## Writing New testcases
- Just create a file ending with `_test.py` in `tests` directory
  - Any subclass of BaseTest with name ending with `Test` will be executed
  - see `run_tests` and `get_testcases` in `core/runner.py` for more info

## Setup
Requires python 3.9 or later.

## Setup for macOS / Linux

### Install tmux
- You need to have tmux installed. See https://github.com/tmux/tmux/wiki

### Python virtual env setup
```
# cd to this directory
cd <path/to/here>
python3 -m venv .venv
.venv/bin/pip install --upgrade pip
.venv/bin/pip install -r requirements.txt
```

### Make sure you build spf
```
# cd to the superfile repo root (parent of this)
cd <superfile_root>
./build.sh
```

### Running testsuite
```
.venv/bin/python3 main.py
```
## Setup for Windows
Coming soon.



### Python virtual env setup
```
# cd to this directory
cd <path/to/here>

# If your python command refers to python3, you can use 'python' below
python3 -m venv .venv
.venv\Scripts\python -m pip install --upgrade pip
.venv\Scripts\pip install -r requirements.txt
```

### Make sure you build spf
```
# cd to the superfile repo root (parent of this)
cd <superfile_root>
go build -o bin/spf.exe
```

### Running testsuite
Notes
- You must keep your focus on the terminal for the entire duration of test run. `pyautogui` sends keypress to process on focus.

```
.venv\Scripts\python main.py
```

## Tips while running tests
- Use `-d` or `--debug` to enable debug logs during test run.
- If you see flakiness in test runs due to superfile being still open, consider using `--close-wait-time` options to increase wait time for superfile to close. Note : For now we have enforcing superfile to close within a specific time window in tests to reduce test flakiness
- Make sure that your hotkeys are set to default hotkeys. Tests use default hotkeys for now.
- Use `-t` or `--tests` to only run specific tests
  - Example `python main.py -d -t RenameTest CopyTest`
- If you see `libtmux` errors like `libtmux.exc.LibTmuxException: ['no server running on /private/tmp/tmux-501/superfile']` Make sure your python version is up to date
