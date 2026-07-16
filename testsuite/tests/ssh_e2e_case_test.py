from core.command_case import GoCommandCaseTest
from core.environment import Environment


class SSHQuickConnectCaseTest(GoCommandCaseTest):
    CASES = ["ssh_quick_connect"]

    def __init__(self, test_env: Environment):
        super().__init__(test_env)

    def command(self) -> list[str]:
        return ["go", "test", "./src/internal", "-run", "^TestSSHQuickConnectCase$", "-count=1"]


class SSHManualConnectCaseTest(GoCommandCaseTest):
    CASES = ["ssh_manual_connect"]

    def __init__(self, test_env: Environment):
        super().__init__(test_env)

    def command(self) -> list[str]:
        return ["go", "test", "./src/internal", "-run", "^TestSSHManualConnectCase$", "-count=1"]


class SSHFailureModesCaseTest(GoCommandCaseTest):
    CASES = ["ssh_failure_modes"]

    def __init__(self, test_env: Environment):
        super().__init__(test_env)

    def command(self) -> list[str]:
        return ["go", "test", "./src/internal", "-run", "^TestSSHFailureModesCase$", "-count=1"]
