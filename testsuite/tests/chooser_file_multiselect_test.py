from pathlib import Path
import time

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst

TESTROOT = Path("chooser_file_multiselect_ops")
DIR1 = TESTROOT / "dir1"
DIR2 = TESTROOT / "dir2"
FILE1 = DIR1 / "file1.txt"
FILE2 = DIR1 / "file2.txt"
CHOOSER_FILE = DIR2 / "chooser_file.txt"


class ChooserFileMultiSelectTest(GenericTestImpl):

    def __init__(self, test_env: Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=DIR1,
            test_dirs=[DIR1, DIR2],
            test_files=[(FILE1, tconst.FILE_TEXT1), (FILE2, tconst.FILE_TEXT1)],
            key_inputs=['v', 'J', 'J', 'e'],
            validate_spf_closed=True,
            close_wait_time=3,
        )

        self.spf_opts += ["--chooser-file", str(self.env.fs_mgr.abspath(CHOOSER_FILE))]

    def end_execution(self) -> None:
        self.logger.debug("Skipping esc key press for chooser file multiselect test")
        time.sleep(self.close_wait_time)
        self.logger.debug("Finished Execution")

    def validate(self) -> bool:
        if not super().validate():
            return False

        try:
            assert self.env.fs_mgr.check_exists(CHOOSER_FILE), f"File {CHOOSER_FILE} does not exists"
            chooser_file_content = self.env.fs_mgr.read_file(CHOOSER_FILE)
            expected = "\n".join([
                str(self.env.fs_mgr.abspath(FILE1)),
                str(self.env.fs_mgr.abspath(FILE2)),
            ])
            assert chooser_file_content == expected, f"Expected '{expected}', got '{chooser_file_content}'"
        except AssertionError as ae:
            self.logger.debug("Test assertion failed : %s", ae, exc_info=True)
            return False

        return True
