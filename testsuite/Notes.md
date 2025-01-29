## Implementation notes

- The `pyautogui` sends input to the process in focus, which is the `spf` subprocess.
- If `spf` is not exited correcly via `q`, it causes wierd vertical tabs in print statements from python
- There is some flakiness in sending of input. Many times, `Ctrl+C` is received as `C` in `spf`
  - If first key is `Ctrl+C`, its always received as `C`
- Note : You must keep your focus on the terminal for the entire duration of test run. `pyautogui` sends keypress to process on focus.

## Pyautogui alternatives
POC with pyautogui as a lot of issues, stated above.

### Linux / MacOS

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

## Windows

- Autohotkey
- ControlSend and SendInput utility in windows
- pywin32 library
  - Create a new GUI window for test
  - Use `win32gui.SendMessage` or `PostMessage`
- References
  - https://www.reddit.com/r/tmux/comments/l580mi/is_there_a_tmuxlike_equivalent_for_windows/