from pathlib import Path
import platform
import uuid
from typing import Optional

from core.base_test import GenericTestImpl
from core.environment import Environment
from core.trash_utils import cleanup_trashed_path, find_trashed_path
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("delete_dir")


class DeleteDirTest(GenericTestImpl):

    def __init__(self, test_env : Environment):
        self.dir1 = TESTROOT / f"spf_trash_dir_{uuid.uuid4().hex}"
        self.nested_dir1 = self.dir1 / "nested1"
        self.nested_dir2 = self.dir1 / "nested2"
        self.file1 = self.nested_dir1 / "file1.txt"
        self.trashed_path : Optional[Path] = None
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[TESTROOT, self.dir1, self.nested_dir1, self.nested_dir2],
            test_files=[(self.file1, tconst.FILE_TEXT1)],
            key_inputs=[keys.KEY_CTRL_D, keys.KEY_ENTER],
            validate_not_exists=[self.dir1, self.nested_dir1, self.nested_dir2, self.file1]
        )

    def validate(self) -> bool:
        if not super().validate():
            return False

        original_path = self.env.fs_mgr.abspath(self.dir1)
        self.trashed_path = find_trashed_path(original_path)
        if self.trashed_path is None:
            self.logger.debug("Could not find deleted directory in trash: %s", original_path)
            return False

        if platform.system() != "Windows":
            nested_file = self.trashed_path / "nested1" / "file1.txt"
            if not nested_file.exists():
                self.logger.debug("Trashed directory did not preserve nested file: %s", nested_file)
                return False
        return True

    def cleanup(self) -> None:
        cleanup_trashed_path(self.trashed_path)
        super().cleanup()
