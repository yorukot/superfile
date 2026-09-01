@aliases.register('spf')
@aliases.uncapturable
@aliases.unthreadable
def _spf(args, stdin=None, stdout=None, stderr=None):
    import os
    import shlex
    import platform
    import subprocess
    from pathlib import Path

    if platform.system() == 'Darwin':
        spf_last_dir = Path(f'{$HOME}/Library/Application Support/superfile/lastdir')
    elif 'XDG_STATE_HOME' in @.env:
        spf_last_dir = Path(f'{$XDG_STATE_HOME}/superfile/lastdir')
    else:
        spf_last_dir = Path(f'{$HOME}/.local/state/superfile/lastdir')

    status_code = subprocess.call(
        ('spf',) + tuple(args), # 'superfile' on nixos
        stdin=stdin,
        stderr=stderr,
        stdout=stdout,
    )

    if status_code == 0 and spf_last_dir.is_file():
        with spf_last_dir.open() as f:
            content = f.read()
            if content:
                _, path = shlex.split(content)
                os.chdir(path)
        spf_last_dir.unlink()

    return status_code

