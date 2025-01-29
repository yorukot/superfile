import win32gui
import win32con
import win32api
import win32process
import subprocess
import time
from ctypes import *

class TerminalController:
    # Virtual Key codes for special keys
    VK_CODES = {
        'ENTER': 0x0D,
        'UP': 0x26,
        'DOWN': 0x28,
        'LEFT': 0x25,
        'RIGHT': 0x27,
        'SPACE': 0x20,
        'TAB': 0x09,
        'ESC': 0x1B
    }
    
    def __init__(self, process_name):
        self.process_name = process_name
        self.hwnd = None
        
    def find_window_by_process(self):
        def callback(hwnd, pid):
            try:
                _, process_id = win32process.GetWindowThreadProcessId(hwnd)
                if process_id == pid:
                    self.hwnd = hwnd
                    return False  # Stop enumeration
            except:
                pass
            return True
            
        # Start the process
        process = subprocess.Popen(self.process_name, 
                                 creationflags=subprocess.CREATE_NEW_CONSOLE)
        time.sleep(1)  # Give process time to start
        
        # Find the window handle
        win32gui.EnumWindows(callback, process.pid)
        return self.hwnd is not None

    def send_key(self, key):
        if not self.hwnd:
            raise Exception("Window handle not found")
            
        if isinstance(key, str):
            if key.upper() in self.VK_CODES:
                vk_code = self.VK_CODES[key.upper()]
            else:
                # For regular characters
                vk_code = win32api.VkKeyScan(key) & 0xFF
                
        # Send key down
        win32api.PostMessage(self.hwnd, win32con.WM_KEYDOWN, vk_code, 0)
        time.sleep(0.05)  # Small delay between down and up
        # Send key up
        win32api.PostMessage(self.hwnd, win32con.WM_KEYUP, vk_code, 0)
        
    def send_keys(self, keys, delay=0.1):
        """Send multiple keys with delay between them"""
        for key in keys:
            self.send_key(key)
            time.sleep(delay)

# Example usage
def main():
    # Initialize controller with your terminal file manager
    controller = TerminalController("superfile")
    
    if controller.find_window_by_process():
        # Example: Navigate using arrow keys and select with Enter
        time.sleep(1)  # Wait for application to be ready
        
        # Send some navigation commands
        commands = ['DOWN', 'DOWN', 'ENTER', 'UP', 'RIGHT']
        controller.send_keys(commands)
    else:
        print("Failed to find terminal window")

if __name__ == "__main__":
    main()