FILE_TEXT1 : str = "This is a sample Text\n"

KEY_DELAY : float       = 0.05 # seconds
OPERATION_DELAY : float = 0.3 # seconds

# 0.3 second was too less for windows
# 0.5 second Github workflow failed for with superfile is still running errors
CLOSE_WAIT_TIME : float     = 1 # seconds
