import libtmux
import time 
import logging
import core.keys as keys
from core.spf_manager import BaseSPFManager

class TmuxSPFManager(BaseSPFManager):
    """
    Tmux based Manager
    After running spf, you can connect to the session via
    tmux -L superfile attach -t spf_session
    Wont work in windows
    """
    # Class variables
    SPF_START_DELAY : float = 0.1 # seconds
    SPF_SOCKET_NAME : str = "superfile"

    # Init should not allocate any resources
    def __init__(self, spf_path : str):
        super().__init__(spf_path)
        
        # Check libtmux version requirement
        min_version = (0, 31, 0)
        current_version_str = libtmux.__version__
        
        # Parse version string to tuple for comparison
        try:
            current_version = tuple(map(int, current_version_str.split('.')[:3]))
        except (ValueError, AttributeError):
            current_version = (0, 0, 0)
        
        if current_version < min_version:
            raise RuntimeError(
                f"libtmux version 0.31.0 or higher is required. "
                f"Current version: {current_version_str}. "
                f"Please upgrade with: pip install 'libtmux>=0.31.0'"
            )
        
        self.logger = logging.getLogger()
        self.server = libtmux.Server(socket_name=TmuxSPFManager.SPF_SOCKET_NAME)
        self.logger.debug("server object : %s", self.server)
        self.spf_session : libtmux.Session = None
        self.spf_pane : libtmux.Pane = None

    def start_spf(self, start_dir : str = None, args : list[str] = None) -> None:
        spf_command = self.spf_path
        if args:
            spf_command += " " + " ".join(args)

        self.logger.debug("windows_command : %s", spf_command)
        

        self.spf_session= self.server.new_session('spf_session',
                window_command=spf_command, 
                start_directory=start_dir)
        time.sleep(TmuxSPFManager.SPF_START_DELAY)
        self.logger.debug("spf_session initialised : %s", self.spf_session)

        # If libtmux version is less than 0.3.1, active_pane does not exist.
        self.spf_pane = self.spf_session.active_pane
        self._is_spf_running = True

    def _send_key(self, key : str) -> None:
        self.logger.debug("sendind key : %s", repr(key))
        self.spf_pane.send_keys(key, enter=False)

    def send_text_input(self, text : str, all_at_once : bool = True) -> None:
        if all_at_once:
            self._send_key(text)
        else:
            for c in text:
                self._send_key(c)

    def send_special_input(self, key : keys.Keys) -> str:
        if key.ascii_code != keys.NO_ASCII:
            self._send_key(chr(key.ascii_code))
        elif isinstance(key, keys.SpecialKeys):
            self._send_key(key.key_name)
        else:
            raise Exception(f"Unknown key : {key}") 
            
    def get_rendered_output(self) -> str:
        return "[Not supported yet]"

    def is_spf_running(self) -> bool:
        self._is_spf_running = (self.spf_session is not None) \
            and (self.spf_session in self.server.sessions)

        return self._is_spf_running

    def close_spf(self) -> None:
        if self.is_spf_running():
            self.server.kill_session(self.spf_session.name)

    # Override
    def runtime_info(self) -> str:
        return str(self.server.sessions)

    def __repr__(self) -> str:
        return f"{self.__class__.__name__}(server : {self.server}, " + \
            f"session : {self.spf_session}, running : {self._is_spf_running})"
