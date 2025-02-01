from abc import ABC

class Keys(ABC):
    def __init__(self, ascii_code : int):
        self.ascii_code = ascii_code
    
    def __repr__(self):
        return f"Key(code={self.ascii_code})"

# Will isinstance of Keys work for object of CtrlKeys ?
class CtrlKeys(Keys):
    def __init__(self, char : str):
        # Only allowing single alphabetic character
        # assert is good here as all objects are defined statically
        assert len(char) == 1
        assert char.isalpha() and char.islower()
        self.char = char
        # Ctrl + A starts at 1
        super().__init__(ord(char) - ord('a') +  1)

# Maybe have keycode
class SpecialKeys(Keys):
    def __init__(self, ascii_code : int, key_name : str):
        super().__init__(ascii_code)
        self.key_name = key_name



KEY_CTRL_A : Keys = CtrlKeys('a')
KEY_CTRL_C : Keys = CtrlKeys('c')
KEY_CTRL_M : Keys = CtrlKeys('m')
KEY_CTRL_R : Keys = CtrlKeys('r')
KEY_CTRL_V : Keys = CtrlKeys('v')

# See https://vimdoc.sourceforge.net/htmldoc/digraph.html#digraph-table for key codes
KEY_BACKSPACE   : Keys = SpecialKeys(8 , "Backspace")
KEY_ENTER       : Keys = SpecialKeys(13, "Enter")
KEY_ESC         : Keys = SpecialKeys(27, "Esc")

