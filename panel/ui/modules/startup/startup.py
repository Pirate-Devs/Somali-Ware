import os


class Startup:
    def __init__(self) -> None:
        self._SomaliWare_logs_dir = os.path.join(
            os.getenv("APPDATA"), "Somali-Ware", "logs"
        )

    def delete_old_logs(self):
        found_logs = []
        for folder_name in os.listdir(self._SomaliWare_logs_dir):
            pass
