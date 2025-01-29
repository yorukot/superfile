## Implementation notes

- The `pyautogui` sends input to the process in focus, which is the `spf` subprocess.
- If `spf` is not exited correcly via `q`, it causes wierd vertical tabs in print statements from python
- There is some flakiness in sending of input. Many times, `Ctrl+C` is received as `C` in `spf`
  - If first key is `Ctrl+C`, its always received as `C`
- Note : You must keep your focus on the terminal for the entire duration of test run. `pyautogui` sends keypress to process on focus.

## Pyautogui alternatives
- POC with pyautogui as a lot of issues, stated above.
- Linux piping
- xdotool
- mkfifo
- tmux
  - Supports full terminal programs
- References
  - https://superuser.com/questions/585398/sending-simulated-keystrokes-in-bash