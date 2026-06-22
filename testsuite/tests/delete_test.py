from pathlib import Path
import uuid
from typing import Optional

from core.base_test import GenericTestImpl
from core.environment import Environment
from core.trash_utils import cleanup_trashed_path, find_trashed_path
import core.test_constants as tconst
import core.keys as keys

TESTROOT = Path("delete_ops")



class DeleteTest(GenericTestImpl):

    def __init__(self, test_env : Environment):
        self.file1 = TESTROOT / f"spf_trash_file_{uuid.uuid4().hex}.txt"
        self.trashed_path : Optional[Path] = None
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[TESTROOT],
            test_files=[(self.file1, tconst.FILE_TEXT1)],
            key_inputs=[keys.KEY_CTRL_D, keys.KEY_ENTER],
            validate_not_exists=[self.file1]
        )

    def validate(self) -> bool:
        if not super().validate():
            return False

        original_path = self.env.fs_mgr.abspath(self.file1)
        self.trashed_path = find_trashed_path(original_path)
        if self.trashed_path is None:
            self.logger.debug("Could not find deleted file in trash: %s", original_path)
            return False
        return True

    def cleanup(self) -> None:
        cleanup_trashed_path(self.trashed_path)
        super().cleanup()


class PermanentlyDeleteTest(GenericTestImpl):

    def __init__(self, test_env : Environment):
        self.file1 = TESTROOT / f"spf_permanent_delete_file_{uuid.uuid4().hex}.txt"
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[TESTROOT],
            test_files=[(self.file1, tconst.FILE_TEXT1)],
            key_inputs=['D', keys.KEY_ENTER],
            validate_not_exists=[self.file1]
        )

    def validate(self) -> bool:
        if not super().validate():
            return False

        original_path = self.env.fs_mgr.abspath(self.file1)
        trashed_path = find_trashed_path(original_path)
        if trashed_path is not None:
            self.logger.debug("Permanently deleted file was found in trash: %s", trashed_path)
            cleanup_trashed_path(trashed_path)
            return False
        return True
