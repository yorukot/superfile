## Implementation notes

- The `pyautogui` sends input to the process in focus, which is the `spf` subprocess.
- If `spf` is not exited correcly via `q`, it causes wierd vertical tabs in print statements from python
- There is some flakiness in sending of input. Many times, `Ctrl+C` is received as `C` in `spf`
  - If first key is `Ctrl+C`, its always received as `C`