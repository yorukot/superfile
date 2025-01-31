import subprocess
import time
import tempfile
import os

def send_keys(process, keys):
    """Send keys to the subprocess."""
    for key in keys:
        print(f"Writing {key}")
        process.stdin.write(key.encode())
        process.stdin.flush()
        time.sleep(0.1)  # Add a slight delay between key presses

def main():
    try:
        with tempfile.TemporaryDirectory() as temp_dir:
            print(f'Temporary directory created at: {temp_dir}')
            
            # Create a sample file
            file1_path = os.path.join(temp_dir, 'example.txt')
            with open(file1_path, 'w') as f:
                f.write("This is a test file.")

            # Start a new command prompt session
            process = subprocess.Popen([r"C:\Users\nitin\Documents\Programming\superfile\bin\spft.exe"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

            print("Command Prompt started.")

            # Send commands to the command prompt
            #send_keys(process, 'echo Hello World\n')
            #send_keys(process, f'type {file1_path}\n')
            send_keys(process, 'exit\n')
            send_keys(process, 'q\n')

            print("done")
            

            # Wait for the process to complete
            process.kill()

            print("Command Prompt session ended.")

    except Exception as e:
        print("Exception during execution:", e)

if __name__ == "__main__":
    main()
