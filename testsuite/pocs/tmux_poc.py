import os
import time
import tempfile
import libtmux
import pathlib

def spf_input(p : libtmux.Pane) -> None:
    p.send_keys(chr(3), enter=False)
    #p.send_keys("C-c", enter=False)
    p.send_keys(chr(22), enter=False)
    #p.send_keys("C-v", enter=False)
    
    # Maybe this is bad for efficiency, hence async is great.
    time.sleep(0.1)
    p.send_keys('q', enter=False)

def main() -> None:
    try:
        with tempfile.TemporaryDirectory() as temp_dir:
            print(f'Temporary directory created at: {temp_dir}')
            # division operator is overloaded
            dir1:pathlib.Path = pathlib.Path(temp_dir) / "dir1"
            file1_path:pathlib.Path = dir1 / "file1.txt"
            file1_cpy_path:pathlib.Path = dir1 / "file1(1).txt"

            # setup
            os.makedirs(dir1, exist_ok=True)
            with open(file1_path, 'w') as f:
                f.write("This is a test file.")

            server:libtmux.Server = libtmux.Server()
            
            # We can have two levels of logging in test, verbose and non verbose
            print(f"Tmux server started : {server}")

            spf_session:libtmux.Session = server.new_session('spf_session',
                window_command='/Users/kuknitin/Workspace/kuknitin/superfile/bin/spf', 
                start_directory=dir1)
            print(f"Tmux session started : {spf_session}")
            
            
            spf_input(spf_session.active_pane)

            time.sleep(0.1)

            if os.path.isfile(file1_cpy_path):
                print("File copied successfully!")
            else:
                print("File copy failed.")

            # Might be a little inefficient, as it will iterate through all sessions
            if server.sessions.count(spf_session) == 0:
                print("spf session has exited")
            else:
                raise Exception("spf session still lingering around")
            
    except Exception as e:
        print("Exception during test : ", e)

main()

