from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys
import time

TESTROOT = Path("search_mode_ops")
DIR1 = TESTROOT / "dir1"
NESTED = TESTROOT / "nested" / "deep"
FILE1 = DIR1 / "main.go"
FILE2 = NESTED / "main.go"
OTHER = TESTROOT / "other.txt"


class SearchModeTest(GenericTestImpl):
    """Recursive search mode: open via Z, match relative paths, open results.

    Navigation is validated through file-system effects: file creation happens
    in the panel's current directory, and copy+paste creates a suffixed copy of
    the focused item.
    """
    def __init__(self, test_env: Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=TESTROOT,
            test_dirs=[TESTROOT, DIR1, NESTED],
            test_files=[
                (FILE1, tconst.FILE_TEXT1),
                (FILE2, tconst.FILE_TEXT1),
                (OTHER, tconst.FILE_TEXT1),
            ],
        )

    def search_and_wait(self, query: str) -> None:
        # Z opens search mode; the query types into the searchbar. Wait for
        # debounce + walk + streaming before interacting with the results.
        self.env.spf_mgr.send_text_input("Z")
        time.sleep(tconst.KEY_DELAY)
        self.env.spf_mgr.send_text_input(query)
        time.sleep(1.0)

    def create_file(self, name: str) -> None:
        # ctrl+n opens the create prompt. The first keypress is consumed by
        # the model, so send a backspace no-op before the file name.
        self.env.spf_mgr.send_special_input(keys.KEY_CTRL_N)
        time.sleep(tconst.KEY_DELAY)
        self.env.spf_mgr.send_special_input(keys.KEY_BACKSPACE)
        time.sleep(tconst.KEY_DELAY)
        self.env.spf_mgr.send_text_input(name)
        time.sleep(tconst.KEY_DELAY)
        self.env.spf_mgr.send_special_input(keys.KEY_ENTER)
        time.sleep(tconst.OPERATION_DELAY)

    # Override
    def test_execute(self) -> None:
        self.start_spf()
        time.sleep(tconst.OPERATION_DELAY)
        # No-op keypress: the first keypress after startup is not registered
        self.env.spf_mgr.send_text_input("a")
        time.sleep(tconst.KEY_DELAY)

        # Opening a file result navigates to its containing directory and
        # focuses the file: copy+paste of the focused main.go yields a
        # suffixed copy inside dir1.
        self.search_and_wait("main")
        self.env.spf_mgr.send_special_input(keys.KEY_ENTER)
        time.sleep(tconst.KEY_DELAY)
        self.env.spf_mgr.send_special_input(keys.KEY_CTRL_C)
        time.sleep(tconst.KEY_DELAY)
        self.env.spf_mgr.send_special_input(keys.KEY_PASTE)
        time.sleep(tconst.OPERATION_DELAY)
        assert self.env.fs_mgr.check_exists(DIR1 / "main(1).go"), "copied main.go should be in dir1"

        # Opening a directory result makes it the new focused directory: a
        # file created afterwards lands inside it.
        self.env.spf_mgr.send_text_input("h")  # Back to the search root
        time.sleep(tconst.KEY_DELAY)
        self.search_and_wait("deep")
        self.env.spf_mgr.send_special_input(keys.KEY_ENTER)
        time.sleep(tconst.KEY_DELAY)
        self.create_file("created_deep.txt")
        assert self.env.fs_mgr.check_exists(NESTED / "created_deep.txt"), "file should be created in the opened directory"

        # Cancelling with Esc leaves the panel on the search root: a file
        # created afterwards lands at the root.
        self.env.spf_mgr.send_text_input("h")
        time.sleep(tconst.KEY_DELAY)
        self.env.spf_mgr.send_text_input("h")
        time.sleep(tconst.KEY_DELAY)
        self.search_and_wait("other")
        self.env.spf_mgr.send_special_input(keys.KEY_ESC)
        time.sleep(tconst.KEY_DELAY)
        self.create_file("created_root.txt")
        assert self.env.fs_mgr.check_exists(TESTROOT / "created_root.txt"), "file should be created at the search root after cancel"

        self.end_execution()
