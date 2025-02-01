from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("cut_ops")
DIR1 = TESTROOT / "dir1"
DIR2 = TESTROOT / "dir2"
FILE1 = DIR1 / "file1.txt"
FILE1_CUT1 = DIR2 / "file1.txt"



class CutTest(GenericTestImpl):

    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=DIR1,
            test_dirs=[DIR1, DIR2],
            test_files=[(FILE1, tconst.FILE_TEXT1)],
            key_inputs=[keys.KEY_CTRL_X, keys.KEY_LEFT, keys.KEY_DOWN, 
                keys.KEY_ENTER, keys.KEY_CTRL_V],
            validation_files=[FILE1_CUT1]
        )
