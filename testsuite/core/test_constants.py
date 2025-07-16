import platform
from pathlib import Path
FILE_TEXT1 : str = "This is a sample Text\n"

KEY_DELAY : float       = 0.05 # seconds
OPERATION_DELAY : float = 0.3 # seconds

# 0.3 second was too less for windows
# 0.5 second Github workflow failed for with superfile is still running errors
START_WAIT_TIME : float     = 0.5 # seconds
CLOSE_WAIT_TIME : float     = 0.5 # seconds

# Platform specific consts
FILE_CREATE_COMMAND : str   = "touch"
if platform.system() == "Windows" :
    FILE_CREATE_COMMAND = "ni"

CONF_DIR = Path(__file__).parent.parent.parent / "src" / "superfile_config"

CONFIG_FILE = CONF_DIR / "config.toml"
HOTKEY_FILE = CONF_DIR / "hotkeys.toml"
