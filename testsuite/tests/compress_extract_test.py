from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("ce_ops")
DIR1 = TESTROOT / "dir1"
FILE1 = DIR1 / "file1"
FILE2 = DIR1 / "file2"

DIR1_ZIPPED = TESTROOT / "dir1.zip"

DIR1_EXTRACTED = TESTROOT / "dir1(1)" / "dir1"
FILE1_EXTRACTED = DIR1_EXTRACTED / "file1"
FILE2_EXTRACTED = DIR1_EXTRACTED / "file2"


class CompressExtractTest(GenericTestImpl):
    """Test compression and extraction

    Args:
        GenericTestImpl (_type_): _description_
    """
    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[DIR1],
            test_files=[(FILE1, tconst.FILE_TEXT1), (FILE2, tconst.FILE_TEXT1)],
            key_inputs=[keys.KEY_CTRL_A, keys.KEY_DOWN, keys.KEY_CTRL_E],
            validate_exists=[DIR1, DIR1_ZIPPED, DIR1_EXTRACTED, FILE1_EXTRACTED, FILE2_EXTRACTED]
        )
