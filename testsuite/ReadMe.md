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