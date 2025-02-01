import logging
from abc import ABC, abstractmethod
from core.environment import Environment


class BaseTest(ABC):
    """Base class for all tests
    The idea is to have independency among each test.
    And for each test to have full control on its environment, execution, and validation.
    """
    def __init__(self, test_env : Environment):
        self.test_env = test_env
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
