import core.keys as keys
from core.spf_manager import BaseSPFManager


class NullSPFManager(BaseSPFManager):
    def start_spf(self, start_dir: str = None, args: list[str] = None) -> None:
        raise RuntimeError("SPF manager is unavailable for this test environment")

    def send_text_input(self, text: str, all_at_once: bool = False) -> None:
        raise RuntimeError("SPF manager is unavailable for this test environment")

    def send_special_input(self, key: keys.Keys) -> None:
        raise RuntimeError("SPF manager is unavailable for this test environment")

    def get_rendered_output(self) -> str:
        return "[SPF manager unavailable]"

    def is_spf_running(self) -> bool:
        return False

    def close_spf(self) -> None:
        self._is_spf_running = False
