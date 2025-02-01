import argparse
import logging
import sys

from core.fs_manager import TestFSManager

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
    try:
        t = TestFSManager()
        logger.info(t)

        t.makedirs('1/2/3')
        t.create_file("1/2/3/1.txt")
        logger.info(t.tree('1'))
        input("Press enter to exit ...")
        t.cleanup() 
    except Exception as e:
        logger.error("Exception while running tests : {%s}", e)
    finally:
        pass

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