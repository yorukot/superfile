import subprocess
import win32gui
import win32con
import win32api
import time

def send_keys_to_window(hwnd, text):
    # Bring window to front
    win32gui.ShowWindow(hwnd, win32con.SW_NORMAL)
    win32gui.SetForegroundWindow(hwnd)
    time.sleep(0.5)  # Wait for window to be ready
    
    # Send each character
    for char in text:
        vk_code = win32api.VkKeyScan(char) & 0xFF
        win32api.PostMessage(hwnd, win32con.WM_CHAR, ord(char), 0)
        time.sleep(0.1)  # Small delay between chars
    
    # Send Enter key
    win32api.PostMessage(hwnd, win32con.WM_KEYDOWN, win32con.VK_RETURN, 0)
    win32api.PostMessage(hwnd, win32con.WM_KEYUP, win32con.VK_RETURN, 0)

# Start Notepad
notepad = subprocess.Popen(['notepad.exe'])
time.sleep(1)  # Wait for Notepad to start

# Find Notepad window
hwnd = win32gui.FindWindow("Notepad", None)

# Send text if window found
if hwnd:
    send_keys_to_window(hwnd, "hello")
else:
    print("Notepad window not found")