from pathlib import Path

from core.base_test import GenericTestImpl
from core.environment import Environment
import core.test_constants as tconst
import core.keys as keys
import time

TESTROOT = Path("empty_panel_ops")
DIR1 = TESTROOT / "dir1"


class EmptyPanelTest(GenericTestImpl):
    """
    Validate that spf doesn't crashes when we try to 
    perform operations on empty file panel
    """
    def __init__(self, test_env : Environment):
        super().__init__(
            test_env=test_env,
            test_root=TESTROOT,
            start_dir=DIR1,
            test_dirs=[DIR1],
            key_inputs=[ 
                keys.KEY_CTRL_C,    # Try copy
                keys.KEY_CTRL_X,    # Try cut
                keys.KEY_CTRL_D,    # Try delete
                keys.KEY_PASTE,     # Try paste
                keys.KEY_CTRL_R,    # Try rename
                keys.KEY_CTRL_P,    # Try copy location
                'e',                # Try open with editor
                keys.KEY_ENTER,
                keys.KEY_RIGHT,
                keys.KEY_CTRL_A,    # Try archiving
                keys.KEY_CTRL_E,    # Try extract
                'v',                # Try going to Select mode
                'J',                # Try select down  
                'K',                # Try select up
                'A',                # select all
                'v',
                '.',                # Try toggle dotfiles                 
                ],
            # Makes sure spf doesn't crashes
            validate_spf_running=True
        )

    # Override
    def test_execute(self) -> None:
        self.start_spf()
        self.send_input()    
        time.sleep(tconst.OPERATION_DELAY)
        # Intentionally not closing spf to ensure it remains running,
        # which is verified by the validate_spf_running flag which is set 
        # to true for this testcase
    
    
        
    
