import logging
import time
from abc import ABC, abstractmethod
from core.environment import Environment
from pathlib import Path
from typing import Union, List, Tuple
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
    def setup(self) -> None:
        """Set up the required things for test
        """

    @abstractmethod
    def test_execute(self) -> None:
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
        test_dirs : List[Path],
        test_files : List[Tuple[Path, str]],
        key_inputs : List[Union[keys.Keys,str]],
        validate_exists : List[Path] = [],
        validate_not_exists : List[Path] = []):
        super().__init__(test_env)
        self.test_root = test_root
        self.start_dir = start_dir
        self.test_dirs = test_dirs
        self.test_files = test_files
        self.key_inputs = key_inputs
        self.validate_exists = validate_exists
        self.validate_not_exists = validate_not_exists
    
    def setup(self) -> None:
        for dir_path in self.test_dirs:
            self.env.fs_mgr.makedirs(dir_path)
        for file_path, data in self.test_files:
            self.env.fs_mgr.create_file(file_path, data)
        
        self.logger.debug("Current file structure : \n%s",
            self.env.fs_mgr.tree(self.test_root))


    def test_execute(self) -> None:
        """Execute the test
        """
        # Start in DIR1
        self.env.spf_mgr.start_spf(self.env.fs_mgr.abspath(self.start_dir))

        assert self.env.spf_mgr.is_spf_running(), "Superfile is not running"

        for cur_input in self.key_inputs:
            if isinstance(cur_input, keys.Keys):
                self.env.spf_mgr.send_special_input(cur_input)
            else:
                assert isinstance(cur_input, str), "Invalid input type"
                self.env.spf_mgr.send_text_input(cur_input)
            time.sleep(tconst.KEY_DELAY)

        time.sleep(tconst.OPERATION_DELAY)
        self.env.spf_mgr.send_special_input(keys.KEY_ESC)    
        time.sleep(tconst.CLOSE_WAIT_TIME)
        self.logger.debug("Finished Execution")

    def validate(self) -> bool:
        """Validate that test passed. Log exception if failed.
        Returns:
            bool: True if validation passed
        """
        self.logger.debug("spf_manager info : %s, Current file structure : \n%s",
            self.env.spf_mgr.runtime_info(), self.env.fs_mgr.tree(self.test_root))
        try:
            assert not self.env.spf_mgr.is_spf_running(), "Superfile is still running"
            for file_path in self.validate_exists:
                assert self.env.fs_mgr.check_exists(file_path), f"File {file_path} does not exists"
            
            for file_path in self.validate_not_exists:
                assert not self.env.fs_mgr.check_exists(file_path), f"File {file_path} exists" 
        except AssertionError as ae:
            self.logger.debug("Test assertion failed : %s", ae, exc_info=True)
            return False
                
        return True
    
    def __repr__(self) -> str:
        return f"{self.__class__.__name__}"

