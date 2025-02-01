import libtmux
import time 
from abc import ABC, abstractmethod

from core.keys import Keys


class BaseSPFManager(ABC):

    def __init__(self, spf_path : str):
        self.spf_path = spf_path
        # _ denotes the internal variables, anyone should not directly read/modify
        self._is_spf_running : bool = False

    @abstractmethod
    def start_spf(self, start_dir : str = None) -> None:
        pass 
    
    @abstractmethod
    def send_text_input(self, text : str, all_at_once : bool = False) -> None:
        pass 

    @abstractmethod
    def send_special_input(self, key : Keys) -> None:
        pass 

    @abstractmethod
    def get_rendered_output(self) -> str:
        pass
    
    
    @abstractmethod
    def is_spf_running(self) -> bool:
        """
        We allow using _is_spf_running variable for efficiency
        But this method should give the true state, although this might have some calculations
        """
        return self._is_spf_running
    
    @abstractmethod
    def close_spf(self) -> None:
        """
        Close spf if its running and cleanup any other resources
        """


class TmuxSPFManager(BaseSPFManager):
    """
    Tmux based Manager
    """
    # Class variables
    SPF_START_DELAY = 0.1 # seconds
    SPF_SOCKET_NAME = "superfile"

    # Init should not allocate any resources
    def __init__(self, spf_path : str):
        super().__init__(spf_path)
        self.server = libtmux.Server(socket_name=TmuxSPFManager.SPF_SOCKET_NAME)
        self.spf_session : libtmux.Session = None
        self.spf_pane : libtmux.Pane = None

    def start_spf(self, start_dir : str = None) -> None:
        self.spf_session= self.server.new_session('spf_session',
                window_command=self.spf_path, 
                start_directory=start_dir)
        time.sleep(TmuxSPFManager.SPF_START_DELAY)

        self.spf_pane = self.spf_session.active_pane
        self._is_spf_running = True

    def _send_key(self, key : str) -> None:
        self.spf_pane.send_keys(key, enter=False)

    def send_text_input(self, text : str, all_at_once : bool = True) -> None:
        if all_at_once:
            self._send_key(text)
        else:
            for c in text:
                self._send_key(c)

    def send_special_input(self, key : Keys) -> str:
        self._send_key(chr(key.ascii_code))

    def get_rendered_output(self) -> str:
        return "[Not supported yet]"

    def is_spf_running(self) -> bool:
        self._is_spf_running = (
            (self.spf_session != None) 
            and (self.server.sessions.count(self.spf_session) == 1))

        return self._is_spf_running

    def close_spf(self) -> None:
        if self.is_spf_running():
            self.server.kill_session(self.spf_session.name)

    def __repr__(self) -> str:
        return f"{self.__class__.__name__}(server : {self.server}, " + \
            f"session : {self.spf_session}, running : {self._is_spf_running})"


class PyAutoGuiSPFManager(BaseSPFManager):
    
    def __init__(self, spf_path : str):
        super().__init__(spf_path)
        self.spf_process = None


    def start_spf(self, start_dir : str = None) -> None:
        pass 
    
    def send_text_input(self, text : str, all_at_once : bool = False) -> None:
        pass 

    def send_special_input(self, key : Keys) -> None:
        pass 

    def get_rendered_output(self) -> str:
        pass
    
    
    def is_spf_running(self) -> bool:
        """
        We allow using _is_spf_running variable for efficiency
        But this method should give the true state, although this might have some calculations
        """
        return self._is_spf_running
    
    def close_spf(self) -> None:
        """
        Close spf if its running and cleanup any other resources
        """
