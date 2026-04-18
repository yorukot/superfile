from pathlib import Path
import time

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("save_file_ops")
DIR1 = TESTROOT / "dir1"
SAVE_OUT = TESTROOT / "save_out.txt"


class SaveFileTest(GenericTestImpl):

    def __init__(self, test_env: Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=DIR1,
            test_dirs=[DIR1],
            key_inputs=[keys.KEY_CTRL_R, "download.txt", keys.KEY_ENTER, "E"],
            validate_spf_closed=True,
            close_wait_time=3,
        )

        self.spf_opts += [
            "--save-file",
            str(self.env.fs_mgr.abspath(SAVE_OUT)),
            str(self.env.fs_mgr.abspath(DIR1)),
        ]

    def end_execution(self) -> None:
        self.logger.debug("Skipping esc key press for save file test")
        time.sleep(self.close_wait_time)
        self.logger.debug("Finished Execution")

    def validate(self) -> bool:
        if not super().validate():
            return False

        try:
            assert self.env.fs_mgr.check_exists(SAVE_OUT), f"File {SAVE_OUT} does not exists"
            assert self.env.fs_mgr.check_exists(DIR1 / "download.txt"), "download placeholder was not created at target path"
            save_output = self.env.fs_mgr.read_file(SAVE_OUT)
            expected = str(self.env.fs_mgr.abspath(DIR1 / "download.txt"))
            assert save_output == expected, f"Expected '{expected}', got '{save_output}'"
        except AssertionError as ae:
            self.logger.debug("Test assertion failed : %s", ae, exc_info=True)
            return False

        return True
