import pyperclip
# ------ Clipboard utils

# This creates a layer of abstraction.
# Now the user of the fuction doesn't need to import pyperclip
# or need to even know what pyperclip was used.
def get_sys_clipboard_text() -> str :
    return pyperclip.paste()