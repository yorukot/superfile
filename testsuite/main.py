import argparse
import logging
import sys
from pathlib import Path

from core.runner import run_tests
import core.test_constants as tconst


def configure_logging(debug : bool = False) -> None:
    # Prefer stdout instead of default stderr
    handler = logging.StreamHandler(sys.stdout)
    
    # 7s to align all log levelnames - WARNING is the largest level, with size 7
    handler.setFormatter(logging.Formatter(
        '[%(asctime)s - %(levelname)7s] %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S'
    ))


    logger = logging.getLogger()
    logger.addHandler(handler)

    if debug:
        logger.setLevel(logging.DEBUG)
    else:
        logger.setLevel(logging.INFO)

    logging.getLogger("libtmux").setLevel(logging.WARNING)

def main():
    # Setup argument parser
    parser = argparse.ArgumentParser(description='superfile testsuite')
    parser.add_argument('-d', '--debug',action='store_true',
                        help='Enable debug logging')
    parser.add_argument('--close-wait-time', type=float,
                        help='Override default wait time after closing spf')
    parser.add_argument('--spf-path', type=str,
                        help='Override the default spf executable path(../bin/spf) under test')
    
    # Parse arguments
    args = parser.parse_args()
    if args.close_wait_time is not None:
        tconst.CLOSE_WAIT_TIME = args.close_wait_time
    
    configure_logging(args.debug)
        
    # Default path
    # We maybe should run this only in main.py file.
    spf_path = Path(__file__).parent.parent / "bin" / "spf"

    if args.spf_path is not None:
        spf_path = Path(args.spf_path)
    # Resolve any symlinks, and make it absolute
    spf_path = spf_path.resolve()

    run_tests(spf_path)


main()
