import time 
import subprocess
import pyautogui
import core.keys as keys
from core.spf_manager import BaseSPFManager

class PyAutoGuiSPFManager(BaseSPFManager):
    """Manage SPF via subprocesses and pyautogui
    Cross platform, but it globally takes over the input, so you need the terminal 
    constantly on focus during test run
    """
    SPF_START_DELAY : float = 0.5
    def __init__(self, spf_path : str):
        super().__init__(spf_path)
        self.spf_process = None


    def start_spf(self, start_dir : str = None, args : list[str] = None) -> None:
        spf_args = [self.spf_path]
        if args :
            spf_args += args
        spf_args.append(start_dir)

        self.spf_process = subprocess.Popen(spf_args,
            stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        time.sleep(PyAutoGuiSPFManager.SPF_START_DELAY)

        # Need to send a sample keypress otherwise it ignores first keypress
        self.send_text_input('x')
        
    
    def send_text_input(self, text : str, all_at_once : bool = False) -> None:
        if all_at_once :
            pyautogui.write(text)
        else:
            for c in text:
                pyautogui.write(c)

    def send_special_input(self, key : keys.Keys) -> None:
        if isinstance(key, keys.CtrlKeys):
            pyautogui.hotkey('ctrl', key.char)
        elif isinstance(key, keys.SpecialKeys):
            pyautogui.press(key.key_name.lower())
        else:
            raise Exception(f"Unknown key : {key}") 

    def get_rendered_output(self) -> str:
        return "[Not supported yet]" 
    
    
    def is_spf_running(self) -> bool:
        self._is_spf_running = (self.spf_process is not None) and (self.spf_process.poll() is None)
        return self._is_spf_running
    
    def close_spf(self) -> None:
        if self.spf_process is not None:
            self.spf_process.terminate()
    
    # Override
    def runtime_info(self) -> str:
        if self.spf_process is None:
            return "[No process]"
        else:
            return f"[PID : {self.spf_process.pid}, poll : {self.spf_process.poll()}]"  



