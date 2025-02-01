import os
import subprocess
import time
import pyautogui
import tempfile

def spf_input():
    pyautogui.hotkey('x') # Random key press. First key cant be Ctrl+C for some reason
    pyautogui.hotkey('ctrl', 'c')  # Simulate Ctrl+C
    time.sleep(0.5)  # Wait for copy action to complete
    pyautogui.hotkey('ctrl', 'v')  # Simulate Ctrl+V
    time.sleep(0.5)
    pyautogui.hotkey('q')
    time.sleep(0.5)

def main():
    process = None # To make it accessible at all scopes
    try:
        with tempfile.TemporaryDirectory() as temp_dir:
            print(f'Temporary directory created at: {temp_dir}')
            
            dir1 = os.path.join(temp_dir, "dir1")
            file1_path = os.path.join(dir1, "file1.txt")
            file1_cpy_path = os.path.join(dir1, "file1(1).txt")
            spf_out_filepath = os.path.join(temp_dir, "out.txt")
            spf_err_filepath = os.path.join(temp_dir, "err.txt")
            # setup
            os.makedirs(dir1, exist_ok=True)
            with open(file1_path, 'w') as f:
                f.write("This is a test file.")

            with open(spf_out_filepath, 'w') as fout, open(spf_err_filepath, 'w') as ferr:
                # Start Superfile in a subprocess
                process = subprocess.Popen(
                    ['../bin/spf', dir1],
                    #stdout=fout, 
                    #stderr=ferr,
                    )
                # Wait for it to load.
                time.sleep(1)
                spf_input()
                

                if os.path.isfile(file1_cpy_path):
                    print("File copied successfully!")
                else:
                    print("File copy failed.")

                if process.poll() is not None:
                    print("spf process exited successfuly")
                else:
                    print("spf process is still running")
    except Exception as e:
        print("Exception during test : ", e)
    finally:
        if process is not None:
            process.terminate()

main()

