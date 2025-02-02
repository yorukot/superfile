# Overview
This is to document the behviour of tmux, and how could we use it in testsuite

# Tmux concepts and working info
- Tmux creates a main server process, and one new process for each session.
<img width="1147" alt="image" src="https://github.com/user-attachments/assets/c701a56f-49ef-43a2-9b76-8f613f5e219a" />
- `-s` and `-n` for window naming.
- We have prefix keys to send commands to tmux. 

# Sample usage with spf

## Sending keys to termux and controlling from outside.
<img width="1147" alt="image" src="https://github.com/user-attachments/assets/8ca67e0c-ece6-465c-84b1-755c8fa0a67e" />
<img width="1147" alt="image" src="https://github.com/user-attachments/assets/e4aeb0b9-e184-41a0-89eb-9174082eb884" />

# Knowledge sharing
- `tmux new 'spf'` - Run spf in tmux
- `tmux attach -t <session_name>` attach to an existing session. You can have two windows duplicating same behaviour.
- `tmux kill-session -t <session_name>` kill session
- `Ctrl+B`+`:` - Enter commands
- `Ctrl+B`+`D` - Detach from session
- `:source ~/.tmux.conf` - Change the config of running server
- We have already a wrapper libary for termux in python !!!!!
- How to send key press/tmux commands to the process ?


# References
- https://github.com/tmux/tmux/wiki/Getting-Started
- https://tao-of-tmux.readthedocs.io/en/latest/manuscript/10-scripting.html#controlling-tmux-send-keys
- https://github.com/tmux-python/libtmux
