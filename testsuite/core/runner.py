
from core.spf_manager import BaseSPFManager, TmuxSPFManager, PyAutoGuiSPFManager
from core.fs_manager import TestFSManager
from core.environment import Environment
from core.base_test import BaseTest

import logging
import platform
import importlib
from pathlib import Path


logger = logging.getLogger()

def get_testcases(test_env : Environment) -> list[BaseTest]:
    res : list[BaseTest] = []
    test_dir = Path(__file__).parent.parent / "tests"
    for test_file in test_dir.glob("*_test.py"):
        # Import dynamically
        module_name = test_file.stem 
        module = importlib.import_module(f"tests.{module_name}")
        for attr_name in dir(module):
            attr = getattr(module, attr_name)
            if isinstance(attr, type) and attr is not BaseTest and issubclass(attr, BaseTest) \
                and  attr_name.endswith("Test"):
                logger.debug("Found a testcase %s, in module %s", attr_name, module_name)
                res.append(attr(test_env))
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
        testcases = get_testcases(test_env)
        logger.debug("Testcases : %s", testcases)
        for t in testcases:
            t.setup()
            t.test_execute()
            cnt_executed += 1

            if t.validate():
                cnt_passed += 1
            elif stop_on_fail:
                break
        
        logger.info("Finised running %s test. %s passed", cnt_executed, cnt_passed)
    finally:
        # Make sure of cleanup
        # This is still not full proof, as if what happens when TestFSManager __init__ fails ?
        test_env.cleanup()


