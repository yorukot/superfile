import logging
from spf_manager import BaseSPFManager, TmuxSPFManager
from superfile.testsuite.core.base_test import BaseExecutor, CopyTest

def run_tests(spf_path : str, stop_on_fail : bool = True) -> None:
    spf_manager = TmuxSPFManager(spf_path)
    tests : list[BaseExecutor] = [CopyTest(spf_manager, fs)]

    cnt_passed : int = 0
    cnt_executed : int = 0
    for t in tests:
        t.setup()
        t.test_execute()
        cnt_executed += 1
        
        if t.validate():
            cnt_passed += 1
        elif stop_on_fail:
            break


