import subprocess
import time
import win32gui
import win32con
import win32api

# Open Notepad
subprocess.Popen("notepad.exe")

# Wait for Notepad to open
time.sleep(2)

# Find Notepad window
hwnd = win32gui.FindWindow(None, "Untitled - Notepad")
hwnd = win32gui.FindWindowEx(hwnd, None, "Edit", None)
print(hwnd)

# Send text "hellow" and press Enter
for char in "hellow":
    win32api.SendMessage(hwnd, win32con.WM_CHAR, ord(char), 0)
win32api.SendMessage(hwnd, win32con.WM_KEYDOWN, win32con.VK_RETURN, 0)
win32api.SendMessage(hwnd, win32con.WM_KEYUP, win32con.VK_RETURN, 0)
