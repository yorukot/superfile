from core.spf_manager import BaseSPFManager
from core.fs_manager import TestFSManager
from core.environment import Environment
from core.base_test import BaseTest

import logging
import platform
import importlib
from pathlib import Path
from typing import List


# Preferred importing at the top level
if platform.system() == "Windows" :
    # Conditional import is needed to make it work on linux
    # importing pyautogui on linux can cause errors.
    from core.pyautogui_manager import PyAutoGuiSPFManager
else:
    from core.tmux_manager import TmuxSPFManager

logger = logging.getLogger()

def get_testcases(test_env : Environment, only_run_tests : List[str] = None) -> List[BaseTest]:
    res : List[BaseTest] = []
    test_dir = Path(__file__).parent.parent / "tests"
    for test_file in test_dir.glob("*_test.py"):
        # Import dynamically
        module_name = test_file.stem 
        module = importlib.import_module(f"tests.{module_name}")
        for attr_name in dir(module):
            if only_run_tests is not None and attr_name not in only_run_tests:
                continue
            attr = getattr(module, attr_name)
            if isinstance(attr, type) and attr is not BaseTest and issubclass(attr, BaseTest) \
                and  attr_name.endswith("Test"):
                logger.debug("Found a testcase %s, in module %s", attr_name, module_name)
                res.append(attr(test_env))
    return res

def run_tests(spf_path : Path, stop_on_fail : bool = True, only_run_tests : List[str] = None) -> bool:
    """Runs tests

    Args:
        spf_path (Path): Path of spf binary under test
        stop_on_fail (bool, optional): Whether to stop on failures. Defaults to True.
        only_run_tests (List[str], optional): Only specific test to run. Defaults to None.

    Returns:
        bool: Whether run was successful
    """    
    # is this str conversion needed ?

    spf_manager : BaseSPFManager = None 
    if platform.system() == "Windows" :
        spf_manager = PyAutoGuiSPFManager(str(spf_path))
    else:
        spf_manager = TmuxSPFManager(str(spf_path))
        
    fs_manager = TestFSManager()

    test_env = Environment(spf_manager, fs_manager)
    cnt_passed : int = 0
    cnt_executed : int = 0
    try:    
        testcases : List[BaseTest] = get_testcases(test_env, only_run_tests=only_run_tests)
        logger.info("Testcases : %s", testcases)
        for t in testcases:
            logger.info("Running test %s", t)
            t.setup()
            t.test_execute()
            cnt_executed += 1
            passed : bool = t.validate()
            t.cleanup()

            if passed:
                logger.info("Passed test %s", t)
                cnt_passed += 1
            else:
                logger.error("Failed test %s", t)
                if stop_on_fail:
                    break
        
        logger.info("Finished running %s test. %s passed", cnt_executed, cnt_passed)
    finally:
        # Make sure of cleanup
        # This is still not full proof, as if what happens when TestFSManager __init__ fails ?
        test_env.cleanup()

    return cnt_passed == cnt_executed
        


