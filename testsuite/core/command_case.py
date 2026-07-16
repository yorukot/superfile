import os
import subprocess
from abc import abstractmethod
from pathlib import Path
from typing import Dict, List

from core.base_test import BaseTest
from core.environment import Environment


class GoCommandCaseTest(BaseTest):
    WORKDIR = Path(__file__).resolve().parents[2]

    def __init__(self, test_env: Environment):
        super().__init__(test_env)
        self.stdout: str = ""
        self.stderr: str = ""
        self.returncode: int = -1

    @classmethod
    def requires_spf_manager(cls) -> bool:
        return False

    @abstractmethod
    def command(self) -> List[str]:
        pass

    def environment(self) -> Dict[str, str]:
        env = os.environ.copy()
        env["SSH_AUTH_SOCK"] = ""
        return env

    def setup(self) -> None:
        return None

    def test_execute(self) -> None:
        cmd = self.command()
        self.logger.info("Running command-backed SSH case: %s", " ".join(cmd))
        result = subprocess.run(
            cmd,
            cwd=self.WORKDIR,
            env=self.environment(),
            capture_output=True,
            text=True,
            check=False,
        )
        self.returncode = result.returncode
        self.stdout = result.stdout
        self.stderr = result.stderr

    def validate(self) -> bool:
        if self.stdout:
            self.logger.info("stdout for %s:\n%s", self, self.stdout)
        if self.stderr:
            self.logger.info("stderr for %s:\n%s", self, self.stderr)
        return self.returncode == 0

    def cleanup(self) -> None:
        return None
