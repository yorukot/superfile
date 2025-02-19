from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
from core.utils import get_sys_clipboard_text
import core.test_constants as tconst
import core.keys as keys
from assertpy import assert_that
import time

TESTROOT = Path("nav_ops")
DIR1 = TESTROOT / "dir1"
FILE1 = TESTROOT / "file1"
FILE2 = TESTROOT / "file2"

# Temporarily disabled, till we fix xclip does not works in github actions 
class NavCopyPathTest_Disabled(GenericTestImpl):
    """Test navigation, and Copying of path
    """
    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[TESTROOT, DIR1],
            test_files=[(FILE1, tconst.FILE_TEXT1), (FILE2, tconst.FILE_TEXT1)]
        )
    
    # Override
    def test_execute(self) -> None:
        self.start_spf()
        time.sleep(tconst.OPERATION_DELAY)
        # > dir1 
        #   file1
        #   file2

        self.env.spf_mgr.send_special_input(keys.KEY_CTRL_P)
        time.sleep(tconst.KEY_DELAY)
        assert_that(get_sys_clipboard_text()).is_equal_to(str(self.env.fs_mgr.abspath(DIR1)))

        self.env.spf_mgr.send_special_input(keys.KEY_DOWN)
        time.sleep(tconst.KEY_DELAY)
        #   dir1 
        # > file1
        #   file2
        self.env.spf_mgr.send_special_input(keys.KEY_CTRL_P)
        time.sleep(tconst.KEY_DELAY)
        assert_that(get_sys_clipboard_text()).is_equal_to(str(self.env.fs_mgr.abspath(FILE1)))

        self.env.spf_mgr.send_special_input(keys.KEY_DOWN)
        time.sleep(tconst.KEY_DELAY)
        #   dir1 
        #   file1
        # > file2
        self.env.spf_mgr.send_special_input(keys.KEY_CTRL_P)
        time.sleep(tconst.KEY_DELAY)
        assert_that(get_sys_clipboard_text()).is_equal_to(str(self.env.fs_mgr.abspath(FILE2)))

        self.env.spf_mgr.send_special_input(keys.KEY_UP)
        time.sleep(tconst.KEY_DELAY)
        #   dir1 
        # > file1
        #   file2
        self.env.spf_mgr.send_special_input(keys.KEY_CTRL_P)
        time.sleep(tconst.KEY_DELAY)
        assert_that(get_sys_clipboard_text()).is_equal_to(str(self.env.fs_mgr.abspath(FILE1)))

        self.env.spf_mgr.send_special_input(keys.KEY_DOWN)
        time.sleep(tconst.KEY_DELAY)
        self.env.spf_mgr.send_special_input(keys.KEY_DOWN)
        time.sleep(tconst.KEY_DELAY)
        # > dir1 
        #   file1
        #   file2
        self.env.spf_mgr.send_special_input(keys.KEY_CTRL_P)
        time.sleep(tconst.KEY_DELAY)
        assert_that(get_sys_clipboard_text()).is_equal_to(str(self.env.fs_mgr.abspath(DIR1)))
        
        self.end_execution()
