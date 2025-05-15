import logging
from  tempfile import TemporaryDirectory
from pathlib import Path
import os
from io import StringIO

class TestFSManager:
    """Manage the temporary files for test and the cleanup
    """    
    def __init__(self):
        self.logger = logging.getLogger()
        self.logger.debug("Initialized %s", self.__class__.__name__) 
        self.temp_dir_obj = TemporaryDirectory()
        self.temp_dir = Path(self.temp_dir_obj.name)
    
    def abspath(self, relative_path : Path) -> Path:
        return self.temp_dir / relative_path
    
    def check_exists(self, relative_path : Path) -> bool:
        return self.abspath(relative_path).exists()
    
    def read_file(self, relative_path: Path) -> str:
        content = ""
        try:
            with open(self.temp_dir / relative_path, 'r', encoding="utf-8") as f:
                content = f.read()
        except FileNotFoundError:
            self.logger.error("File not found: %s", relative_path)
        except PermissionError:
            self.logger.error("Permission denied when reading file: %s", relative_path)
        return content

    def makedirs(self, relative_path : Path) -> None:
        # Overloaded '/' operator
        os.makedirs(self.temp_dir / relative_path, exist_ok=True)
    
    def create_file(self, relative_path : Path, data : str = "") -> None:
        """Create files
        Make sure directories exist
        Args:
            relative_path (Path): Relative path from test root
        """
        with open(self.temp_dir / relative_path, 'w', encoding="utf-8") as f:
            f.write(data)

    def tree(self, relative_root : Path = None) -> str:
        if relative_root is None:
            root = self.temp_dir
        else:
            root = self.temp_dir / relative_root
        res = StringIO()
        for item in root.rglob('*'):
            path_str = str(item.relative_to(root))
            if item.is_dir():
                res.write(f"D-{path_str}\n")
            else:
                res.write(f"F-{path_str}\n")
        return res.getvalue()

    def cleanup(self) -> None:
        """Cleaup the temporary directory
        Its okay to forget it though, it will be cleaned on program exit then.
        """
        self.temp_dir_obj.cleanup()

    def __repr__(self) -> str:
        return f"{self.__class__.__name__}(temp_dir = {self.temp_dir})"
