from spf_manager import BaseSPFManager
from fs_manager import TestFSManager

class Environment:
    def __init__(self, spf_manager : BaseSPFManager, fs_manager : TestFSManager ):
        self.spf_manager = spf_manager
        self.fs_manager = fs_manager

    def cleanup(self):
        self.spf_manager.close_spf()
        self.fs_manager.cleanup()