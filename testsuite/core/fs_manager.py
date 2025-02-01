import fs
import fs.copy
import fs.memoryfs
import logging
import io 

# Relevant constants
FILE_TEXT1 : str = "This is a sample Text\n"


# Abstract the usage of in memory fs
class TestFSManager:

    # Class variables for efficiency
    TEST_1_COPY_ROOT_DIR = "op1/copy_op"

    def __init__(self):
        # Gotta close it too
        # usage of memoryfs to quickly replicate the whole test filesystem anywhere
        self.memfs = fs.memoryfs.MemoryFS()
        self.logger = logging.getLogger()

        self.setup_memfs()
        self.logger.debug("Initialized TestDirManager") 
        tree_out = io.StringIO()
        self.memfs.tree(file=tree_out)
        
        # Using lazy logging -  W1201
        self.logger.debug("Directory structure : \n%s", tree_out.getvalue())
    
    def copy_to_dir(self, dst_dir : str) -> None:
        fs.copy.copy_fs(self.memfs, dst_dir)
    
    def writef_memfs(self, filepath : str, content : str = "") -> None:
        with self.memfs.open(filepath, 'w') as f:
            f.write(content)

    def setup_memfs(self) -> None:
        
        # At root level

        # First group of operations
        self.memfs.makedirs("op1")

        # Copy operation. makedir creates intermediate directories
        self.memfs.makedirs("op1/copy_op/dir1")
        self.memfs.makedirs("op1/copy_op/dir2")
        self.writef_memfs("op1/copy_op/dir1/file1.txt", FILE_TEXT1)

        # Cut operation
        self.memfs.makedirs("op1/cut_op/dir1")
        self.memfs.makedirs("op1/cut_op/dir2")
        self.writef_memfs("op1/cut_op/dir1/file1.txt", FILE_TEXT1)
        
        # Rename operation
        self.memfs.makedirs("op1/rename_op/dir1")
        self.writef_memfs("op1/rename_op/dir1/file1.txt", FILE_TEXT1)

    def cleanup(self) -> None:
        self.memfs.close()

