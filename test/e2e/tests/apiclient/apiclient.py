from typing import Any

import requests


class APIClient:
    def __init__(
        self, base_url: str, session_cookie: str, organization_id: str | None = None
    ):
        self.base_url: str = base_url + "/" if not base_url.endswith("/") else base_url
        self.cookies: dict[str, str] = {"storeit_session": session_cookie}
        self.headers = {"X-Organization-Id": organization_id} if organization_id else {}

    def set_organization(self, organization_id: str):
        self.headers = {"X-Organization-Id": organization_id}

    def get(self, path: str) -> requests.Response:
        return requests.get(
            self.base_url + path, cookies=self.cookies, headers=self.headers
        )

    def post(self, path: str, data: dict[str, Any]) -> requests.Response:
        return requests.post(
            self.base_url + path, json=data, cookies=self.cookies, headers=self.headers
        )

    def put(self, path: str, data: dict[str, Any]) -> requests.Response:
        return requests.put(
            self.base_url + path, json=data, cookies=self.cookies, headers=self.headers
        )

    def patch(self, path: str, data: dict[str, Any]) -> requests.Response:
        return requests.patch(
            self.base_url + path, json=data, cookies=self.cookies, headers=self.headers
        )

    def delete(self, path: str) -> requests.Response:
        return requests.delete(
            self.base_url + path, cookies=self.cookies, headers=self.headers
        )
