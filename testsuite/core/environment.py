from core.spf_manager import BaseSPFManager
from core.fs_manager import TestFSManager

class Environment:
    """Manage test environment
    Manage cleanup of environment and other stuff at a single place
    """    
    def __init__(self, spf_manager : BaseSPFManager, fs_manager : TestFSManager ):
        self.spf_mgr = spf_manager
        self.fs_mgr = fs_manager

    def cleanup(self):
        self.spf_mgr.close_spf()
        self.fs_mgr.cleanup()