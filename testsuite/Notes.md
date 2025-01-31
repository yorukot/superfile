# Implementation notes

- The `pyautogui` sends input to the process in focus, which is the `spf` subprocess.
- If `spf` is not exited correcly via `q`, it causes wierd vertical tabs in print statements from python
- There is some flakiness in sending of input. Many times, `Ctrl+C` is received as `C` in `spf`
  - If first key is `Ctrl+C`, its always received as `C`
- Note : You must keep your focus on the terminal for the entire duration of test run. `pyautogui` sends keypress to process on focus.

## Input to spf

### Pyautogui alternatives
POC with pyautogui as a lot of issues, stated above.

#### Linux / MacOS

- xdotool
  - Seems complicated. It wont be able to manage spf process that well
- mkfifo / Manual linux piping
  - Too much manual work to send inputs, even if it works
- tmux
  - Supports full terminal programs and has a python wrapper library
  - See `docs/tmux.md`
  - Not available for windows
- References
  - https://superuser.com/questions/585398/sending-simulated-keystrokes-in-bash

#### Windows

- Autohotkey
  - No better than pyautogui
- ControlSend and SendInput utility in windows
  - Isn't that just for C# / C++ code ?
- Python ctypes
  - https://stackoverflow.com/questions/62189991/how-to-wrap-the-sendinput-function-to-python-using-ctypes
- pywin32 library
  - Create a new GUI window for test
  - Use `win32gui.SendMessage` or `win32gui.PostMessage`
  - Probably the correct way, but I havent been able to get it working.
  - First we need to get it send input to a sample window like notepad, etc. Then we can make superfile work
- pywinpty
  - Heavy installations requirements. Needs Rust, and Visual studio build tools.
  - Rust cargo not found
    - Needs rust 
  - link.exe not found (` the msvc targets depend on the msvc linker but link.exe was not found` )
    - Needs to install Visual Studio Build Tools (build tools and spectre mitigated libs)
    - Had to manually find link.exe and put it on the PATH
  - You might get error of unable to find mspdbcore.dll (I havent been able to solve it so far)
    - https://stackoverflow.com/questions/67328795/c1356-unable-to-find-mspdbcore-dll
- References
  - https://www.reddit.com/r/tmux/comments/l580mi/is_there_a_tmuxlike_equivalent_for_windows/

## Directory setup
- Programmatic setup is better.
- We could keep test directory setup as a config file - json/yaml/toml 
- or as a hardcoded python dict

## Tests and Validation
- Each tests starts independently, so there is no strict order
- Hardcoded validations . Predefined test, where each test has start dir, key press, and validations
- We could have a base Class test. where check(), input(), init(), methods would be overrided
- It allows greater flexibility in terms of testcases.
- Abstraction layer for spf init, teardown and inputm