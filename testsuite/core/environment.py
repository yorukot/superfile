from core.spf_manager import BaseSPFManager
from core.fs_manager import TestFSManager

class Environment:
    """Manage test environment
    Manage cleanup of environment and other stuff at a single place
    """    
    def __init__(self, spf_manager : BaseSPFManager, fs_manager : TestFSManager ):
        self.spf_manager = spf_manager
        self.fs_manager = fs_manager

    def cleanup(self):
        self.spf_manager.close_spf()
        self.fs_manager.cleanup()