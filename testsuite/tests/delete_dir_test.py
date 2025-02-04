from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("delete_dir")
DIR1 = TESTROOT / "dir1"
NESTED_DIR1 = DIR1 / "nested1"
NESTED_DIR2 = DIR1 / "nested2"
FILE1 = NESTED_DIR1 / "file1.txt"


class DeleteDirTest(GenericTestImpl):

    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[TESTROOT, DIR1, NESTED_DIR1, NESTED_DIR2],
            test_files=[(FILE1, tconst.FILE_TEXT1)],
            key_inputs=[keys.KEY_CTRL_D, keys.KEY_ENTER],
            validate_not_exists=[DIR1, NESTED_DIR1, NESTED_DIR2, FILE1]
        )
