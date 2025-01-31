import win32gui
import win32con
import win32api
import win32process
import subprocess
import time
from ctypes import *

import win32gui
import win32con
import win32api
import time

class SimpleWindow:
    def __init__(self):
        # Register the Window class
        window_class = win32gui.WNDCLASS()
        window_class.lpszClassName = "SimpleWindowClass"
        window_class.lpfnWndProc = {
            win32con.WM_DESTROY: self.on_destroy,
            win32con.WM_PAINT: self.on_paint,
        }
        
        # Register the class
        self.class_atom = win32gui.RegisterClass(window_class)
        
        # Create the window
        self.hwnd = win32gui.CreateWindow(
            self.class_atom,                   # Class name
            "My Initial Title",                # Window title
            win32con.WS_OVERLAPPEDWINDOW,      # Style
            win32con.CW_USEDEFAULT,           # X position
            win32con.CW_USEDEFAULT,           # Y position
            500,                              # Width
            400,                              # Height
            0,                                # Parent
            0,                                # Menu
            0,                                # Instance
            None                              # Additional application data
        )
        
        # Show the window
        win32gui.ShowWindow(self.hwnd, win32con.SW_SHOW)
        win32gui.UpdateWindow(self.hwnd)
    
    def on_destroy(self, hwnd, message, wparam, lparam):
        """Called when window is closed"""
        win32gui.PostQuitMessage(0)
        return True
    
    def on_paint(self, hwnd, message, wparam, lparam):
        """Handle window painting"""
        paint_struct = win32gui.PAINTSTRUCT()
        hdc = win32gui.BeginPaint(hwnd, paint_struct)
        
        # Paint white background
        rect = win32gui.GetClientRect(hwnd)
        win32gui.FillRect(hdc, rect, win32gui.GetStockObject(win32con.WHITE_BRUSH))
        
        # Add some text
        win32gui.SetTextColor(hdc, win32api.RGB(0, 0, 0))
        win32gui.SetBkMode(hdc, win32con.TRANSPARENT)
        win32gui.DrawText(
            hdc, 
            "Hello, Win32!", 
            -1, 
            rect, 
            win32con.DT_SINGLELINE | win32con.DT_CENTER | win32con.DT_VCENTER
        )
        
        win32gui.EndPaint(hwnd, paint_struct)
        return 0
    
    def change_title(self, new_title):
        """Change the window title"""
        win32gui.SetWindowText(self.hwnd, new_title)
    
    def message_loop(self):
        """Run the message loop"""
        while True:
            try:
                msg = win32gui.GetMessage(None, 0, 0)
                win32gui.TranslateMessage(msg)
                win32gui.DispatchMessage(msg)
            except:
                break

def main_test():
    window = SimpleWindow()
    
    # Change title after 2 seconds
    def change_title_delayed():
        time.sleep(2)
        window.change_title("New Window Title!")
        time.sleep(2)
    
    # Start a separate thread to change the title
    import threading
    threading.Thread(target=change_title_delayed, daemon=True).start()
    
    # Run the message loop
    window.message_loop()


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
        self.process = None
        print("Init done")
        
    def find_window_by_process(self):
        def callback(hwnd, pid):
            #print(f"windosws found a window with pid {pid}, and its handle info : {hwnd}")
            try:
                arg2, process_id = win32process.GetWindowThreadProcessId(hwnd)
                print(f"arg2 : {arg2}, process_id : {process_id}")
                if process_id == pid:
                    self.hwnd = hwnd
                    return False  # Stop enumeration
            except Exception as e:
                print("Exception e ", e)
            return True
            
        # Start the process
        self.process = subprocess.Popen(self.process_name, 
                                 creationflags=subprocess.CREATE_NEW_CONSOLE)
        time.sleep(1)  # Give process time to start
        print(f"started spf with pid={self.process.pid}")
        # Find the window handle
        win32gui.EnumWindows(callback, self.process.pid)
        return self.hwnd is not None
    
    def test(self):
        
        win32gui.MoveWindow(self.hwnd, 100, 100, 100, 100, True)
        time.sleep(0.2)
        win32gui.SetWindowText(self.hwnd, "My new window title")
        time.sleep(0.2)
        win32gui.SetWindowText(self.hwnd, "My old window title")
        
        time.sleep(0.2)
        win32gui.SetWindowText(self.hwnd, "My new window title")
        time.sleep(0.2)
        win32gui.SetWindowText(self.hwnd, "My old window title")
        win32gui.MoveWindow(self.hwnd, 100, 100, 100, 100, True)

    def send_key(self, key):
        win32gui.SetForegroundWindow(self.hwnd)
        if not self.hwnd:
            raise Exception("Window handle not found")
            
        if isinstance(key, str):
            if key.upper() in self.VK_CODES:
                vk_code = self.VK_CODES[key.upper()]
            else:
                # For regular characters
                vk_code = win32api.VkKeyScan(key) & 0xFF
                
        # Send key down
        win32gui.PostMessage(self.hwnd, win32con.WM_KEYDOWN, vk_code, 0)
        
        time.sleep(0.05)  # Small delay between down and up
        # Send key up
        win32gui.PostMessage(self.hwnd, win32con.WM_KEYUP, vk_code, 0)
        
    def send_keys(self, keys, delay=0.1):
        """Send multiple keys with delay between them"""
        for key in keys:
            self.send_key(key)
            time.sleep(delay)

def window_info(hwnd):
    """Get detailed information about a window"""
    if not hwnd:
        return None
        
    info = {
        'title': win32gui.GetWindowText(hwnd),
        'class_name': win32gui.GetClassName(hwnd),
        'rect': win32gui.GetWindowRect(hwnd),
        'visible': win32gui.IsWindowVisible(hwnd),
        'enabled': win32gui.IsWindowEnabled(hwnd),
        'foreground': (hwnd == win32gui.GetForegroundWindow())
    }
    
    # Get process ID
    try:
        _, pid = win32process.GetWindowThreadProcessId(hwnd)
        info['process_id'] = pid
    except:
        info['process_id'] = None
        
    return info

def print_info(hwnd):
    print("\nActive Window Information:")
    active_info = window_info(hwnd)
    for key, value in active_info.items():
        print(f"{key}: {value}")

# Example usage
def main():
    #main_test()
    print_info(win32gui.GetForegroundWindow())
    
    
    # Initialize controller with your terminal file manager
    controller = TerminalController("spft")
    
    if controller.find_window_by_process():
        # Example: Navigate using arrow keys and select with Enter
        time.sleep(1)  # Wait for application to be ready

        print_info(controller.hwnd)

        controller.test()

        print_info(controller.hwnd)

        
        # Send some navigation commands
        commands = ['DOWN', 'DOWN', 'ENTER', 'UP', 'RIGHT']
        controller.send_keys(commands)
        win32gui.SendMessage(controller.hwnd, win32con.WM_KEYDOWN, win32con.VK_RETURN, 0)
        win32gui.SendMessage(controller.hwnd, win32con.WM_KEYDOWN, win32con.VK_RETURN, 0)
        win32gui.SendMessage(controller.hwnd, win32con.WM_KEYDOWN, win32con.VK_RETURN, 0)

    else:
        print("Failed to find terminal window")
    
    time.sleep(2)
    controller.process.terminate()

if __name__ == "__main__":
    main()