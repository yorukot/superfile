from pathlib import Path
import time

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst

TESTROOT = Path("chooser_file_ops")
DIR1 = TESTROOT / "dir1"
DIR2 = TESTROOT / "dir2"
FILE1 = DIR1 / "file1.txt"
CHOOSER_FILE = DIR2 / "chooser_file.txt"



class ChooserFileTest(GenericTestImpl):

    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=DIR1,
            test_dirs=[DIR1, DIR2],
            test_files=[(FILE1, tconst.FILE_TEXT1)],
            key_inputs=['e'],
            validate_spf_closed=True,
            close_wait_time=3
        )
    
    # Override
    def start_spf(self) -> None:
        self.env.spf_mgr.start_spf(self.env.fs_mgr.abspath(self.start_dir), 
            ["--chooser-file", str(self.env.fs_mgr.abspath(CHOOSER_FILE))])
        assert self.env.spf_mgr.is_spf_running(), "Superfile is not running"

    # Override
    def end_execution(self) -> None:
        self.logger.debug("Skipping esc key press for Chooser file test")
        time.sleep(self.close_wait_time)
        self.logger.debug("Finished Execution")
    # Override
    def validate(self) -> bool:
        if not super().validate():
            return False
        
        try:
            assert self.env.fs_mgr.check_exists(CHOOSER_FILE), f"File {CHOOSER_FILE} does not exists"
            chooser_file_content = self.env.fs_mgr.read_file(CHOOSER_FILE)
            assert chooser_file_content == str(self.env.fs_mgr.abspath(FILE1)), \
                f"Expected '{self.env.fs_mgr.abspath(FILE1)}', got '{chooser_file_content}'"

        except AssertionError as ae:
            self.logger.debug("Test assertion failed : %s", ae, exc_info=True)
            return False
                
        return True