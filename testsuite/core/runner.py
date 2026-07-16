from core.spf_manager import BaseSPFManager
from core.fs_manager import TestFSManager
from core.environment import Environment
from core.base_test import BaseTest
from core.null_spf_manager import NullSPFManager

import logging
import platform
import importlib
import inspect
import shutil
from pathlib import Path
from typing import List, Optional, Type


# Preferred importing at the top level
if platform.system() == "Windows" :
    # Conditional import is needed to make it work on linux
    # importing pyautogui on linux can cause errors.
    from core.pyautogui_manager import PyAutoGuiSPFManager
else:
    from core.tmux_manager import TmuxSPFManager

logger = logging.getLogger()

CASE_FILE_HINTS = {
    "ssh_quick_connect": ["ssh_e2e_case_test.py"],
    "ssh_manual_connect": ["ssh_e2e_case_test.py"],
    "ssh_failure_modes": ["ssh_e2e_case_test.py"],
}

def _matches_case_filter(test_class: Type[BaseTest], only_run_cases: Optional[List[str]]) -> bool:
    if only_run_cases is None:
        return True
    case_names = getattr(test_class, "CASES", [])
    return any(case_name in case_names for case_name in only_run_cases)

def _candidate_test_files(test_dir: Path, only_run_cases: Optional[List[str]]) -> List[Path]:
    if only_run_cases is None:
        return sorted(test_dir.glob("*_test.py"))

    hinted_files = set()
    for case_name in only_run_cases:
        hinted_files.update(CASE_FILE_HINTS.get(case_name, []))

    if not hinted_files:
        return sorted(test_dir.glob("*_test.py"))

    return [test_dir / file_name for file_name in sorted(hinted_files)]

def get_testcase_classes(only_run_tests : List[str] = None, only_run_cases: List[str] = None) -> List[Type[BaseTest]]:
    res : List[Type[BaseTest]] = []
    test_dir = Path(__file__).parent.parent / "tests"
    for test_file in _candidate_test_files(test_dir, only_run_cases):
        # Import dynamically
        module_name = test_file.stem 
        module = importlib.import_module(f"tests.{module_name}")
        for attr_name in dir(module):
            if only_run_tests is not None and attr_name not in only_run_tests:
                continue
            attr = getattr(module, attr_name)
            if not isinstance(attr, type):
                continue
            if not _matches_case_filter(attr, only_run_cases):
                continue
            if attr is not BaseTest and issubclass(attr, BaseTest) \
                and not inspect.isabstract(attr) \
                and  attr_name.endswith("Test"):
                logger.debug("Found a testcase %s, in module %s", attr_name, module_name)
                res.append(attr)
    return res

def get_testcases(test_env : Environment, only_run_tests : List[str] = None,
                  only_run_cases: List[str] = None) -> List[BaseTest]:
    testcase_classes = get_testcase_classes(only_run_tests=only_run_tests, only_run_cases=only_run_cases)
    return [testcase_class(test_env) for testcase_class in testcase_classes]

def _tmux_skip_reason() -> Optional[str]:
    if shutil.which("tmux") is None:
        return "tmux executable is not available"
    return None

def _build_spf_manager(spf_path: Path, testcase_classes: List[Type[BaseTest]]) -> tuple[BaseSPFManager, Optional[str]]:
    requires_spf = any(test_class.requires_spf_manager() for test_class in testcase_classes)
    if not requires_spf:
        return NullSPFManager(str(spf_path)), None

    if not spf_path.exists():
        return NullSPFManager(str(spf_path)), f"spf binary not found at {spf_path}"

    if platform.system() == "Windows" :
        return PyAutoGuiSPFManager(str(spf_path)), None

    skip_reason = _tmux_skip_reason()
    if skip_reason is not None:
        return NullSPFManager(str(spf_path)), skip_reason

    try:
        return TmuxSPFManager(str(spf_path)), None
    except Exception as exc:
        return NullSPFManager(str(spf_path)), str(exc)

def run_tests(spf_path : Path, stop_on_fail : bool = True, only_run_tests : List[str] = None,
              only_run_cases: List[str] = None) -> bool:
    """Runs tests

    Args:
        spf_path (Path): Path of spf binary under test
        stop_on_fail (bool, optional): Whether to stop on failures. Defaults to True.
        only_run_tests (List[str], optional): Only specific test to run. Defaults to None.

    Returns:
        bool: Whether run was successful
    """
    # is this str conversion needed ?

    testcase_classes = get_testcase_classes(only_run_tests=only_run_tests, only_run_cases=only_run_cases)
    spf_manager, spf_skip_reason = _build_spf_manager(spf_path, testcase_classes)

    fs_manager = TestFSManager()

    test_env = Environment(spf_manager, fs_manager)
    cnt_passed : int = 0
    cnt_executed : int = 0
    cnt_skipped : int = 0
    try:    
        testcases : List[BaseTest] = [testcase_class(test_env) for testcase_class in testcase_classes]
        logger.info("Testcases : %s", testcases)
        for t in testcases:
            if spf_skip_reason is not None and t.requires_spf_manager():
                logger.warning("Skipped test %s: %s", t, spf_skip_reason)
                cnt_skipped += 1
                continue
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
        
        logger.info("Finished running %s test(s). %s passed, %s skipped", cnt_executed, cnt_passed, cnt_skipped)
    finally:
        if test_env is not None:
            test_env.cleanup()
        elif fs_manager is not None:
            fs_manager.cleanup()

    if cnt_executed == 0:
        logger.error("No tests were executed")
    if cnt_skipped != 0:
        logger.error("Test run is incomplete: %s test(s) were skipped", cnt_skipped)
    return cnt_executed > 0 and cnt_skipped == 0 and cnt_passed == cnt_executed
        
