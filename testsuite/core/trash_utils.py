import logging
import os
import platform
import shutil
import subprocess
import urllib.parse
from pathlib import Path
from typing import Optional


logger = logging.getLogger()


def find_trashed_path(original_path: Path) -> Optional[Path]:
    system = platform.system()
    if system == "Darwin":
        return _find_darwin_trashed_path(original_path)
    if system == "Linux":
        return _find_linux_trashed_path(original_path)
    if system == "Windows":
        return _find_windows_trashed_path(original_path)

    logger.warning("Trash validation is not implemented for platform %s", system)
    return None


def cleanup_trashed_path(trashed_path: Optional[Path]) -> None:
    if trashed_path is None or platform.system() == "Windows":
        return

    try:
        if trashed_path.is_dir() and not trashed_path.is_symlink():
            shutil.rmtree(trashed_path)
        else:
            trashed_path.unlink(missing_ok=True)
    except OSError as err:
        logger.warning("Failed to clean trashed test item %s: %s", trashed_path, err)


def _find_darwin_trashed_path(original_path: Path) -> Optional[Path]:
    trash_dir = Path.home() / ".Trash"
    exact_path = trash_dir / original_path.name
    if exact_path.exists():
        return exact_path

    stem = original_path.stem
    suffix = original_path.suffix
    for candidate in trash_dir.glob(f"{stem}*{suffix}"):
        if candidate.name.startswith(stem):
            return candidate
    return None


def _find_linux_trashed_path(original_path: Path) -> Optional[Path]:
    trash_root = _linux_home_trash_root()
    info_dir = trash_root / "info"
    files_dir = trash_root / "files"
    if not info_dir.exists():
        return None

    expected_path = "Path=" + _escape_trash_info_path(str(original_path))
    for info_file in info_dir.glob("*.trashinfo"):
        try:
            lines = info_file.read_text(encoding="utf-8").splitlines()
        except OSError:
            continue
        if expected_path in lines:
            trash_name = info_file.name[:-len(".trashinfo")]
            trashed_path = files_dir / trash_name
            if trashed_path.exists():
                return trashed_path
    return None


def _linux_home_trash_root() -> Path:
    xdg_data_home = os.environ.get("XDG_DATA_HOME")
    if xdg_data_home and Path(xdg_data_home).is_absolute():
        return Path(xdg_data_home) / "Trash"
    return Path.home() / ".local" / "share" / "Trash"


def _escape_trash_info_path(path: str) -> str:
    return urllib.parse.quote(path, safe="/-_.!~*'()")


def _find_windows_trashed_path(original_path: Path) -> Optional[Path]:
    script = """
$name = $env:SPF_TEST_TRASH_NAME
$shell = New-Object -ComObject Shell.Application
$bin = $shell.Namespace(10)
if ($null -eq $bin) { exit 1 }
$items = @($bin.Items() | Where-Object { $_.Name -eq $name })
if ($items.Count -gt 0) { Write-Output $items[0].Name }
"""
    try:
        env = os.environ.copy()
        env["SPF_TEST_TRASH_NAME"] = original_path.name
        result = subprocess.run(
            ["powershell", "-NoProfile", "-Command", script],
            check=False,
            capture_output=True,
            env=env,
            text=True,
            timeout=10,
        )
    except (OSError, subprocess.SubprocessError) as err:
        logger.warning("Failed to query Windows Recycle Bin: %s", err)
        return None

    if result.returncode == 0 and result.stdout.strip():
        return Path(result.stdout.strip())
    return None
