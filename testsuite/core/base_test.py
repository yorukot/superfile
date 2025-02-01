import logging
import time
from abc import ABC, abstractmethod
from core.environment import Environment
from pathlib import Path
from typing import Union
import core.keys as keys
import core.test_constants as tconst


class BaseTest(ABC):
    """Base class for all tests
    The idea is to have independency among each test.
    And for each test to have full control on its environment, execution, and validation.
    """
    def __init__(self, test_env : Environment):
        self.env = test_env
        self.logger = logging.getLogger()

    @abstractmethod
    def setup(self):
        """Set up the required things for test
        """

    @abstractmethod
    def test_execute(self):
        """Execute the test
        """

    @abstractmethod
    def validate(self) -> bool:
        """Validate that test passed. Log exception if failed.
        Returns:
            bool: True if validation passed
        """

class GenericTestImpl(BaseTest):
    def __init__(self, test_env : Environment,
        test_root : Path,
        start_dir : Path,
        test_dirs : list[Path],
        test_files : list[tuple[Path, str]],
        key_inputs : list[Union[keys.Keys,str]],
        validation_files : list[Path]):
        super().__init__(test_env)
        self.test_root = test_root
        self.start_dir = start_dir
        self.test_dirs = test_dirs
        self.test_files = test_files
        self.key_inputs = key_inputs
        self.validation_files = validation_files
    
    def setup(self):
        for dir_path in self.test_dirs:
            self.env.fs_mgr.makedirs(dir_path)
        for file_path, data in self.test_files:
            self.env.fs_mgr.create_file(file_path, data)
        
        self.logger.debug("Current file structure : \n%s",
            self.env.fs_mgr.tree(self.test_root))


    def test_execute(self):
        """Execute the test
        """
        # Start in DIR1
        self.env.spf_mgr.start_spf(self.env.fs_mgr.abspath(self.start_dir))

        assert self.env.spf_mgr.is_spf_running()

        for cur_input in self.key_inputs:
            if isinstance(cur_input, keys.Keys):
                self.env.spf_mgr.send_special_input(cur_input)
            else:
                assert isinstance(cur_input, str)
                self.env.spf_mgr.send_text_input(cur_input)
            time.sleep(tconst.KEY_DELAY)

        time.sleep(tconst.OPERATION_DELAY)
        self.env.spf_mgr.send_special_input(keys.KEY_ESC)    
        time.sleep(tconst.CLOSE_DELAY)
        self.logger.debug("Finished Execution")

    def validate(self) -> bool:
        """Validate that test passed. Log exception if failed.
        Returns:
            bool: True if validation passed
        """
        self.logger.debug("tmux sessions : %s, Current file structure : \n%s",
            self.env.spf_mgr.server.sessions, self.env.fs_mgr.tree(self.test_root))
        try:
            assert not self.env.spf_mgr.is_spf_running()
            for file_path in self.validation_files:
                assert self.env.fs_mgr.check_exists(file_path)
        except AssertionError as ae:
            self.logger.debug("Test assertion failed : %s", ae)
            return False
                
        return True
    
    def __repr__(self):
        return f"{self.__class__.__name__}"

