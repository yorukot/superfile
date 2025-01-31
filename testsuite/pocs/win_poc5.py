import os
import time
import tempfile
import pathlib
import subprocess
import win32gui
import win32con
import win32api
import winpty
from ctypes import *

class WindowsTerminal:
    def __init__(self):
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
        
        self.Input = Input
        self.Input_I = Input_I
        self.KeyBdInput = KeyBdInput
        
        # Create winpty terminal
        self.agent = winpty.PTY(
            width=80,
            height=25,
            backend=winpty.PTY_BACKEND_WINPTY
        )
        
        # Store the terminal window handle
        time.sleep(1)  # Wait for terminal to start
        self.hwnd = win32gui.FindWindow("ConsoleWindowClass", None)
        
    def send_key(self, key_code, key_down=True):
        """Send a single key event"""
        extra = c_ulong(0)
        ii_ = self.Input_I()
        ii_.ki = self.KeyBdInput(key_code, 0x48, 
                                0 if key_down else 2, 
                                0, pointer(extra))
        x = self.Input(c_ulong(1), ii_)
        windll.user32.SendInput(1, pointer(x), sizeof(x))
        
    def send_keys(self, text):
        """Send a string of text"""
        # Activate window
        if self.hwnd:
            win32gui.ShowWindow(self.hwnd, win32con.SW_NORMAL)
            win32gui.SetForegroundWindow(self.hwnd)
            time.sleep(0.5)
            
        # Type each character
        for char in text:
            vk_code = win32api.VkKeyScan(char) & 0xFF
            self.send_key(vk_code, True)
            time.sleep(0.05)
            self.send_key(vk_code, False)
            time.sleep(0.05)
            
        # Send Enter
        self.send_key(0x0D, True)
        time.sleep(0.05)
        self.send_key(0x0D, False)
        
    def send_ctrl_key(self, key):
        """Send Ctrl+key combination"""
        # Press Ctrl
        self.send_key(win32con.VK_CONTROL, True)
        time.sleep(0.05)
        
        # Press and release the key
        vk_code = win32api.VkKeyScan(key) & 0xFF
        self.send_key(vk_code, True)
        time.sleep(0.05)
        self.send_key(vk_code, False)
        
        # Release Ctrl
        time.sleep(0.05)
        self.send_key(win32con.VK_CONTROL, False)
        
    def close(self):
        """Close the terminal"""
        if self.hwnd:
            win32gui.PostMessage(self.hwnd, win32con.WM_CLOSE, 0, 0)
        self.agent.close()

def main():
    try:
        with tempfile.TemporaryDirectory() as temp_dir:
            print(f'Temporary directory created at: {temp_dir}')
            
            # Create test directory and file
            dir1 = pathlib.Path(temp_dir) / "dir1"
            file1_path = dir1 / "file1.txt"
            file1_cpy_path = dir1 / "file1(1).txt"
            
            os.makedirs(dir1, exist_ok=True)
            with open(file1_path, 'w') as f:
                f.write("This is a test file.")
            
            # Create terminal session
            term = WindowsTerminal()
            print("Terminal session started")
            
            # Change to test directory
            term.send_keys(f'cd "{dir1}"')
            time.sleep(0.5)
            
            # Start your application (replace with your actual command)
            term.send_keys('your_command_here')
            time.sleep(0.5)
            
            # Send Ctrl+C, Ctrl+V, q as in your example
            term.send_ctrl_key('c')
            time.sleep(0.1)
            term.send_ctrl_key('v')
            time.sleep(0.1)
            term.send_keys('q')
            
            time.sleep(0.5)
            
            # Check results
            if os.path.isfile(file1_cpy_path):
                print("File copied successfully!")
            else:
                print("File copy failed.")
            
            # Cleanup
            term.close()
            print("Terminal session closed")
            
    except Exception as e:
        print("Exception during test:", e)

if __name__ == "__main__":
    main()