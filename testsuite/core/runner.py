
from core.spf_manager import BaseSPFManager, TmuxSPFManager, PyAutoGuiSPFManager
from core.fs_manager import TestFSManager
from core.environment import Environment
from core.base_test import BaseTest

import logging
import platform
import importlib
from pathlib import Path


logger = logging.getLogger()

def get_testcases() -> list[BaseTest]:
    res : list[BaseTest] = []
    test_dir = Path(__file__).parent / "tests"
    for test_file in test_dir.glob("*_test.py"):
        # Import dynamically
        module_name = test_file.stem 
        logger.info(test_file)
        module = importlib.import_module(f"core.tests.{module_name}")
        for attr_name in dir(module):
            attr = getattr(module, attr_name)
            if isinstance(attr, type) and attr_name.endswith("Test"):
                logger.debug("Found a testcase %s, in module %s", attr_name, module_name)
                res.append(attr())
    return res


def run_tests(spf_path : Path, stop_on_fail : bool = True) -> None:
    # is this str conversion needed ?

    spf_manager : BaseSPFManager = None 
    if platform.system() == "Windows" :
        spf_manager = PyAutoGuiSPFManager(str(spf_path))
    else:
        spf_manager = TmuxSPFManager(str(spf_path))
        
    fs_manager = TestFSManager()

    test_env = Environment(spf_manager, fs_manager)



    try:
        cnt_passed : int = 0
        cnt_executed : int = 0
        testcases = get_testcases()
        for t in []:
            t.setup()
            t.test_execute()
            cnt_executed += 1

            if t.validate():
                cnt_passed += 1
            elif stop_on_fail:
                break
    finally:
        # Make sure of cleanup
        # This is still not full proof, as if what happens when TestFSManager __init__ fails ?
        test_env.cleanup()


