import subprocess
import win32gui
import win32con
import win32api
import time
from ctypes import *

# Define structures for SendInput
PUL = POINTER(c_ulong)
class KeyBdInput(Structure):
    _fields_ = [("wVk", c_ushort),
                ("wScan", c_ushort),
                ("dwFlags", c_ulong),
                ("time", c_ulong),
                ("dwExtraInfo", PUL)]

class Input_I(Union):
    _fields_ = [("ki", KeyBdInput),]

class Input(Structure):
    _fields_ = [("type", c_ulong),
                ("ii", Input_I)]

def send_key(key_code, key_down=True):
    extra = c_ulong(0)
    ii_ = Input_I()
    ii_.ki = KeyBdInput(key_code, 0x48, 0 if key_down else 2, 0, pointer(extra))
    x = Input(c_ulong(1), ii_)
    windll.user32.SendInput(1, pointer(x), sizeof(x))

def type_string(text):
    for char in text:
        vk_code = win32api.VkKeyScan(char) & 0xFF
        send_key(vk_code, True)
        time.sleep(0.05)
        send_key(vk_code, False)
        time.sleep(0.05)

# Start CMD
cmd = subprocess.Popen(['cmd.exe'], creationflags=subprocess.CREATE_NEW_CONSOLE)
time.sleep(1)  # Wait for CMD to start

# Find CMD window
hwnd = win32gui.FindWindow("ConsoleWindowClass", None)

if hwnd:
    # Activate window
    win32gui.ShowWindow(hwnd, win32con.SW_NORMAL)
    win32gui.SetForegroundWindow(hwnd)
    time.sleep(0.5)  # Wait for window to be active
    
    # Type command
    type_string("echo hello")
    
    # Press Enter
    send_key(0x0D, True)  # Enter key down
    time.sleep(0.05)
    send_key(0x0D, False)  # Enter key up
else:
    print("CMD window not found")