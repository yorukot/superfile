from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("copyw_ops")
FILE1 = TESTROOT / "file1.txt"
FILE1_COPY1 = TESTROOT / "file1(1).txt"

class CopyWTest(GenericTestImpl):
    """Testcase to validate copying with Ctrl+W shortcut 
    """
    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[TESTROOT],
            test_files=[(FILE1, tconst.FILE_TEXT1)],
            key_inputs=[keys.KEY_CTRL_C, keys.KEY_CTRL_W],
            validate_exists=[FILE1, FILE1_COPY1]
        )