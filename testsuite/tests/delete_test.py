from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("delete_ops")
# File1 fails in my mac
FILE1 = TESTROOT / "file_to_delete.txt"



class DeleteTest(GenericTestImpl):

    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[TESTROOT],
            test_files=[(FILE1, tconst.FILE_TEXT1)],
            key_inputs=[keys.KEY_CTRL_D, keys.KEY_ENTER],
            validate_not_exists=[FILE1]
        )
