from pathlib import Path
import time

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("save_file_overwrite_ops")
DIR1 = TESTROOT / "dir1"
FILE1 = DIR1 / "file1.txt"
SAVE_OUT = TESTROOT / "save_out.txt"


class SaveFileOverwriteTest(GenericTestImpl):

    def __init__(self, test_env: Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=DIR1,
            test_dirs=[DIR1],
            test_files=[(FILE1, tconst.FILE_TEXT1)],
            key_inputs=["E", keys.KEY_ENTER],
            validate_spf_closed=True,
            close_wait_time=3,
        )

        self.spf_opts += [
            "--save-file",
            str(self.env.fs_mgr.abspath(SAVE_OUT)),
            str(self.env.fs_mgr.abspath(FILE1)),
        ]

    def end_execution(self) -> None:
        self.logger.debug("Skipping esc key press for save file overwrite test")
        time.sleep(self.close_wait_time)
        self.logger.debug("Finished Execution")

    def validate(self) -> bool:
        if not super().validate():
            return False

        try:
            assert self.env.fs_mgr.check_exists(SAVE_OUT), f"File {SAVE_OUT} does not exists"
            save_output = self.env.fs_mgr.read_file(SAVE_OUT)
            expected = str(self.env.fs_mgr.abspath(FILE1))
            assert save_output == expected, f"Expected '{expected}', got '{save_output}'"
        except AssertionError as ae:
            self.logger.debug("Test assertion failed : %s", ae, exc_info=True)
            return False

        return True
