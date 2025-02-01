import logging
from abc import ABC, abstractmethod
from spf_manager import BaseSPFManager, TmuxSPFManager


class BaseTest(ABC):
    """Base class for all tests
    The idea is to have independency among each test.
    And for each test to have full control on its environment, execution, and validation.
    """    
    def __init__(self, spf_manager : BaseSPFManager, fs):
        self.spf_manager = spf_manager
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
