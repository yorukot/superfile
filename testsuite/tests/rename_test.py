from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("rename_ops")
DIR1 = TESTROOT / "dir1"

# No extension, as in case of extension, the edit cursor appears before the dot, 
# not at the end of filename
FILE1 = DIR1 / "file1"
FILE1_RENAMED = DIR1 / "file2"



class RenameTest(GenericTestImpl):

    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=DIR1,
            test_dirs=[DIR1],
            test_files=[(FILE1, tconst.FILE_TEXT1)],
            key_inputs=[keys.KEY_CTRL_R, keys.KEY_BACKSPACE, '2', keys.KEY_ENTER],
            validation_files=[FILE1_RENAMED]
        )
