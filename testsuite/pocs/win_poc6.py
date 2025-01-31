import subprocess
import time
import win32gui
import win32con
import win32api
import win32process

def send_keys(hwnd, keys):
    """ Sends keystrokes to the given window handle """
    for key in keys:
        win32api.SendMessage(hwnd, win32con.WM_CHAR, ord(key), 0)
        time.sleep(0.05)  # Small delay for reliability

def main():
    # Start a new command prompt process
    process = subprocess.Popen("cmd.exe", creationflags=subprocess.CREATE_NEW_CONSOLE)

    # Wait for it to open
    time.sleep(2)

    # Find the window of the newly created process
    def enum_windows_callback(hwnd, process_id):
        """ Callback function to find window by process ID """
        _, found_pid = win32process.GetWindowThreadProcessId(hwnd)
        if found_pid == process_id:
            global cmd_hwnd
            cmd_hwnd = hwnd

    cmd_hwnd = None
    win32gui.EnumWindows(enum_windows_callback, process.pid)

    if not cmd_hwnd:
        print("Failed to find CMD window.")
        return

    print(f"Found CMD window handle: {cmd_hwnd}")

    # Send some keys (simulate typing commands)
    send_keys(cmd_hwnd, "echo Hello, Windows Terminal!")
    send_keys(cmd_hwnd, "\r")  # Press Enter

    time.sleep(1)

    send_keys(cmd_hwnd, "exit\r")  # Exit CMD

if __name__ == "__main__":
    main()
