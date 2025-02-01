import argparse
import logging
import sys
import tempfile
import time

from ./core/fs_manager import TestFSManager

logger = logging.getLogger()



def configure_logging(debug : bool = False) -> None:
    # Prefer stdout instead of default stderr
    handler = logging.StreamHandler(sys.stdout)
    
    # 7s to align all log levelnames - WARNING is the largest level, with size 7
    handler.setFormatter(logging.Formatter(
        '[%(asctime)s - %(levelname)7s] %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S'
    ))
    
    logger.addHandler(handler)

    if debug:
        logger.setLevel(logging.DEBUG)
    else:
        logger.setLevel(logging.INFO)

def test_main():
    tdm = TestFSManager()
    try:
        run_tests()
        with tempfile.TemporaryDirectory() as temp_dir:
            logger.info(f'Temporary directory created at: {temp_dir}')
            tdm.copy_to_dir(temp_dir)

            time.sleep(10)

            
    except Exception as e:
        logger.error(f"Exception while running tests : {e}")
    finally:
        tdm.cleanup()

def main():
    # Setup argument parser
    parser = argparse.ArgumentParser(description='superfile testsuite')
    parser.add_argument('-d', '--debug', 
                        action='store_true', 
                        help='Enable debug logging')
    
    # Parse arguments
    args = parser.parse_args()

    configure_logging(args.debug)

    test_main()
    

main()