## Coding style rules
- Prefer using strong typing 
- Prefer using type hinting for the first time the variable is declared, and for functions paremeters and return types
- Use `-> None` to explicitly indicate no return value

### Ideas
- Recommended to integrate your IDE with PEP8 to highlight PEP8 violations in real-time
- Enforcing PEP8 via `pylint flake8 pycodestyle` and via pre commit hooks

## Setup 
Requires python 3.9 or later.

## Setup for MacOS / Linux


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
- Note : You must keep your focus on the terminal for the entire duration of test run. `pyautogui` sends keypress to process on focus.
```
.venv/bin/python3 main.py
```
## Setup for Windows
Coming soon.