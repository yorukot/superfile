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
    @abstractmethod
    def cleanup(self) -> None:
        """Any required cleanup after test is done
        """


class GenericTestImpl(BaseTest):
    def __init__(self, test_env : Environment,
        test_root : Path,
        start_dir : Path,
        test_dirs : List[Path],
        key_inputs : List[Union[keys.Keys,str]] = None,
        test_files : List[Tuple[Path, str]] = None,
        validate_exists : List[Path] = None,
        validate_not_exists : List[Path] = None,
        validate_spf_closed: bool = False,
        validate_spf_running: bool = False,
        close_wait_time : float = tconst.CLOSE_WAIT_TIME ):
        super().__init__(test_env)
        self.test_root = test_root
        self.start_dir = start_dir
        self.test_dirs = test_dirs
        self.test_files = test_files
        self.key_inputs = key_inputs
        self.validate_exists = validate_exists
        self.validate_not_exists = validate_not_exists
        self.validate_spf_closed = validate_spf_closed
        self.validate_spf_running = validate_spf_running
        self.close_wait_time = close_wait_time
    
    def setup(self) -> None:
        for dir_path in self.test_dirs:
            self.env.fs_mgr.makedirs(dir_path)
        
        if self.test_files is not None:
            for file_path, data in self.test_files:
                self.env.fs_mgr.create_file(file_path, data)
        
        self.logger.debug("Current file structure : \n%s",
            self.env.fs_mgr.tree(self.test_root))
        
    
    def start_spf(self) -> None:
        self.env.spf_mgr.start_spf(self.env.fs_mgr.abspath(self.start_dir))
        assert self.env.spf_mgr.is_spf_running(), "superfile is not running"

    def end_execution(self) -> None:
        self.env.spf_mgr.send_special_input(keys.KEY_ESC)    
        time.sleep(self.close_wait_time)
        self.logger.debug("Finished Execution")

    def send_input(self) -> None:
        if self.key_inputs is not None:
            for cur_input in self.key_inputs:
                if isinstance(cur_input, keys.Keys):
                    self.env.spf_mgr.send_special_input(cur_input)
                else:
                    assert isinstance(cur_input, str), "Invalid input type"
                    self.env.spf_mgr.send_text_input(cur_input)
                time.sleep(tconst.KEY_DELAY)

    def test_execute(self) -> None:
        """Execute the test
        """
        self.start_spf()
        self.send_input()    
        time.sleep(tconst.OPERATION_DELAY)
        self.end_execution()
        

    def validate(self) -> bool:
        """Validate that test passed. Log exception if failed.
        Returns:
            bool: True if validation passed
        """
        self.logger.debug("spf_manager info : %s, Current file structure : \n%s",
            self.env.spf_mgr.runtime_info(), self.env.fs_mgr.tree(self.test_root))
        try:
            if self.validate_spf_closed :
                assert not self.env.spf_mgr.is_spf_running(), "superfile is still running"
            if self.validate_spf_running :
                assert self.env.spf_mgr.is_spf_running(), "superfile is not running"

            if self.validate_exists is not None:
                for file_path in self.validate_exists:
                    assert self.env.fs_mgr.check_exists(file_path), f"File {file_path} does not exists"

            if self.validate_not_exists is not None:
                for file_path in self.validate_not_exists:
                    assert not self.env.fs_mgr.check_exists(file_path), f"File {file_path} exists" 
        except AssertionError as ae:
            self.logger.debug("Test assertion failed : %s", ae, exc_info=True)
            return False
                
        return True

    def cleanup(self) -> None:
        # Cleanup after test is done
        if self.env.spf_mgr.is_spf_running():
            self.env.spf_mgr.close_spf()
    
    def __repr__(self) -> str:
        return f"{self.__class__.__name__}"

