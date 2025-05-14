import json
import os
import random
import string
import uuid
from typing import Generator

import pytest
import requests

from .apiclient.apiclient import APIClient

API_BASE = os.getenv("API_BASE_URL", "http://localhost:8080")


@pytest.fixture(scope="session")
def test_user() -> dict[str, str]:
    """Create a single test user for all tests"""
    return {
        "email": f"e2e.test.user.{uuid.uuid4()}@example.com",
        "firstName": f"TestUser{random.randint(1, 1000)}",
        "lastName": f"LastName{random.randint(1, 1000)}",
    }


@pytest.fixture(scope="session")
def session_cookie(test_user: dict[str, str]) -> str:
    """Get session cookie for the test user"""
    response = requests.post(f"{API_BASE}/auth/test", json=test_user)
    assert response.status_code == 200, (
        f"Authentication failed: {response.status_code} {response.text}"
    )

    cookies = response.headers.get("Set-Cookie")
    assert cookies, "No cookies received from authentication"

    session_cookie = next(
        (c for c in cookies.split(";") if "storeit_session=" in c), None
    )
    assert session_cookie, "Session cookie not found"

    return session_cookie.split("=")[1].strip()


@pytest.fixture(scope="session")
def api_client(session_cookie: str) -> APIClient:
    client = APIClient(API_BASE, session_cookie)
    response = client.get("/me")
    assert response.status_code == 200
    return client


@pytest.fixture(scope="session")
def api_client_with_organization(
    api_client: APIClient,
    session_cookie: str,
) -> Generator[APIClient, None, None]:
    org_name = str(uuid.uuid4())
    response = api_client.post("/orgs", {"name": org_name, "subdomain": org_name})
    assert response.status_code == 200
    org_data = response.json()["data"]

    client = APIClient(API_BASE, session_cookie)
    response = client.get("/me")
    assert response.status_code == 200
    client.set_organization(org_data["id"])
    yield client
    response = client.delete(f"/orgs/{org_data['id']}")
    assert response.status_code == 204


@pytest.fixture
def organization(api_client: APIClient) -> dict:
    org_name = str(uuid.uuid4())
    response = api_client.post("/orgs", {"name": org_name, "subdomain": org_name})
    assert response.status_code == 200
    data = response.json()["data"]
    assert data["name"] == org_name
    assert data["subdomain"] == org_name
    return data


@pytest.fixture
def organization_unit(api_client_with_organization: APIClient) -> dict:
    unit_name = str(uuid.uuid4())
    alias = generate_random_string()
    address = generate_random_string()

    response = api_client_with_organization.post(
        "/units",
        {
            "name": unit_name,
            "alias": alias,
            "address": address,
        },
    )
    assert response.status_code == 200
    data = response.json()["data"]
    assert data["name"] == unit_name
    assert data["alias"] == alias
    assert data["address"] == address
    return data


@pytest.fixture
def item(api_client_with_organization: APIClient) -> dict:
    name = str(uuid.uuid4())
    description = generate_random_string()
    response = api_client_with_organization.post(
        "/items", {"name": name, "description": description}
    )
    assert response.status_code == 200
    return response.json()["data"]


def generate_random_string(length: int = 4) -> str:
    return "".join(random.choices(string.ascii_letters, k=length))


class TestOrganization:
    def test_full_organization_lifecycle(self, api_client: APIClient) -> None:
        # Create
        org_name = str(uuid.uuid4())
        response = api_client.post("/orgs", {"name": org_name, "subdomain": org_name})
        assert response.status_code == 200, response.text
        org_data = response.json()["data"]
        org_id = org_data["id"]

        # List
        response = api_client.get("/orgs")
        assert response.status_code == 200
        data = response.json()["data"]
        filtered = (x for x in data if x["id"] == org_id)
        assert next(filtered) is not None

        # Get
        api_client.set_organization(org_id)
        response = api_client.get(f"/orgs/{org_id}")
        assert response.status_code == 200, response.text
        data = response.json()["data"]
        assert data["name"] == org_name
        assert data["subdomain"] == org_name

        # Update
        new_name = f"{org_name}_updated"
        response = api_client.put(f"/orgs/{org_id}", {"name": new_name})
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["subdomain"] == org_name

        # Delete
        response = api_client.delete(f"/orgs/{org_id}")
        assert response.status_code == 204


class TestOrganizationUnit:
    def test_full_unit_lifecycle(self, api_client_with_organization: APIClient) -> None:
        # Create
        unit_name = str(uuid.uuid4())
        alias = generate_random_string()
        address = generate_random_string()

        response = api_client_with_organization.post(
            "/units",
            {
                "name": unit_name,
                "alias": alias,
                "address": address,
            },
        )
        assert response.status_code == 200, response.text
        unit_data = response.json()["data"]
        unit_id = unit_data["id"]

        # List
        response = api_client_with_organization.get("/units")
        assert response.status_code == 200
        data = response.json()["data"]
        filtered = (x for x in data if x["alias"] == alias)
        unit = next(filtered)
        assert unit is not None

        # Get
        response = api_client_with_organization.get(f"/units/{unit_id}")
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == unit_name
        assert data["alias"] == alias
        assert data["address"] == address

        # Update
        new_name = f"{unit_name}_updated"
        new_alias = f"{alias}D"
        new_address = f"{address}D"
        response = api_client_with_organization.put(
            f"/units/{unit_id}",
            {"name": new_name, "alias": new_alias, "address": new_address},
        )
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["alias"] == new_alias
        assert data["address"] == new_address

        # Patch
        # new_name = f"{new_name}_patched"
        # response = api_client_with_organization.patch(
        #     f"/units/{unit_id}",
        #     {"name": new_name},
        # )
        # assert response.status_code == 200
        # data = response.json()["data"]
        # assert data["name"] == new_name
        # assert data["alias"] == new_alias
        # assert data["address"] == new_address

        # Delete
        response = api_client_with_organization.delete(f"/units/{unit_id}")
        assert response.status_code == 204, response.text


class TestStorageGroup:
    def test_full_storage_group_lifecycle(
        self, api_client_with_organization: APIClient, organization_unit: dict
    ) -> None:
        # Create
        storage_group_name = str(uuid.uuid4())
        alias = generate_random_string()
        response = api_client_with_organization.post(
            "/storage-groups",
            {
                "name": storage_group_name,
                "alias": alias,
                "unitId": organization_unit["id"],
            },
        )
        assert response.status_code == 200
        storage_group_data = response.json()["data"]
        storage_group_id = storage_group_data["id"]

        # List
        response = api_client_with_organization.get("/storage-groups")
        assert response.status_code == 200
        data = response.json()["data"]
        filtered = (x for x in data if x["alias"] == alias)
        storage_group = next(filtered)
        assert storage_group is not None

        # Get
        response = api_client_with_organization.get(
            f"/storage-groups/{storage_group_id}"
        )
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == storage_group_name
        assert data["alias"] == alias
        assert data["unitId"] == organization_unit["id"]

        # Update
        new_name = f"{storage_group_name}_updated"
        new_alias = f"{alias}D"
        response = api_client_with_organization.put(
            f"/storage-groups/{storage_group_id}",
            {"name": new_name, "alias": new_alias, "unitId": organization_unit["id"]},
        )
        assert response.status_code == 200, storage_group_id
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["alias"] == new_alias
        assert data["unitId"] == organization_unit["id"]

        # # Patch
        # new_name = f"{new_name}_patched"
        # response = api_client_with_organization.patch(
        #     f"/storage-groups/{storage_group_id}",
        #     {"name": new_name},
        # )
        # assert response.status_code == 200
        # data = response.json()["data"]
        # assert data["name"] == new_name
        # assert data["alias"] == new_alias
        # assert data["unitId"] == organization_unit["id"]

        # Delete
        response = api_client_with_organization.delete(
            f"/storage-groups/{storage_group_id}"
        )
        assert response.status_code == 204


class TestCellsGroup:
    def test_full_cells_group_lifecycle(
        self,
        api_client_with_organization: APIClient,
        organization_unit: dict,
    ) -> None:
        # Create
        cells_group_name = str(uuid.uuid4())
        alias = generate_random_string()
        response = api_client_with_organization.post(
            "/cells-groups",
            {
                "name": cells_group_name,
                "alias": alias,
                "unitId": organization_unit["id"],
            },
        )
        assert response.status_code == 200, response.text
        data = response.json()["data"]
        cells_group_id = data["id"]

        # List
        response = api_client_with_organization.get("/cells-groups")
        assert response.status_code == 200
        data = response.json()["data"]
        filtered = (x for x in data if x["alias"] == alias)
        cells_group = next(filtered)
        assert cells_group is not None

        # Get
        response = api_client_with_organization.get(f"/cells-groups/{cells_group_id}")
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == cells_group_name
        assert data["alias"] == alias
        assert data["unitId"] == organization_unit["id"]

        # Update
        new_name = f"{cells_group_name}_updated"
        new_alias = f"{alias}D"
        response = api_client_with_organization.put(
            f"/cells-groups/{cells_group_id}",
            {"name": new_name, "alias": new_alias, "unitId": organization_unit["id"]},
        )
        assert response.status_code == 200, response.text
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["alias"] == new_alias
        assert data["unitId"] == organization_unit["id"]

        # # Patch
        # new_name = f"{new_name}_patched"
        # response = api_client_with_organization.patch(
        #     f"/cells-groups/{cells_group_id}",
        #     {"name": new_name},
        # )
        # assert response.status_code == 200
        # data = response.json()["data"]
        # assert data["name"] == new_name
        # assert data["alias"] == new_alias
        # assert data["unitId"] == organization_unit["id"]

        # Delete
        response = api_client_with_organization.delete(
            f"/cells-groups/{cells_group_id}"
        )
        assert response.status_code == 204


class TestItem:
    def test_full_item_lifecycle(
        self,
        api_client_with_organization: APIClient,
    ) -> None:
        # Create
        name = str(uuid.uuid4())
        description = generate_random_string()

        response = api_client_with_organization.post(
            "/items",
            {"name": name, "description": description},
        )
        assert response.status_code == 200, response.text
        item_data = response.json()["data"]
        item_id = item_data["id"]
        assert item_data["name"] == name
        assert item_data["description"] == description

        # List
        response = api_client_with_organization.get("/items")
        assert response.status_code == 200, response.text
        data = response.json()["data"]
        filtered = (x for x in data if x["id"] == item_id)
        item = next(filtered)
        assert item is not None
        assert item["name"] == name
        assert item["description"] == description

        # Get
        response = api_client_with_organization.get(f"/items/{item_id}")
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == name
        assert data["description"] == description

        # Update
        new_name = f"{name}_updated"
        new_description = f"{description}_updated"
        response = api_client_with_organization.put(
            f"/items/{item_id}",
            {"name": new_name, "description": new_description},
        )
        assert response.status_code == 200, response.text
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["description"] == new_description

        # # Patch
        # new_name = f"{new_name}_patched"
        # response = api_client_with_organization.patch(
        #     f"/items/{item_id}",
        #     {"name": new_name},
        # )
        # assert response.status_code == 200
        # data = response.json()["data"]
        # assert data["name"] == new_name
        # assert data["description"] == new_description

        # Delete
        response = api_client_with_organization.delete(f"/items/{item_id}")
        assert response.status_code == 204

    def test_item_variant_lifecycle(
        self, item: dict, api_client_with_organization: APIClient
    ) -> None:
        item_id = item["id"]
        variant_name = str(uuid.uuid4())
        variant_description = generate_random_string()

        response = api_client_with_organization.post(
            f"/items/{item_id}/variants",
            {"name": variant_name},
        )
        assert response.status_code == 200, response.text
        variant_data = response.json()["data"]
        variant_id = variant_data["id"]

        # List
        response = api_client_with_organization.get(f"/items/{item_id}/variants")
        assert response.status_code == 200
        data = response.json()["data"]
        assert len(data) == 1
        assert data[0]["id"] == variant_id
        assert data[0]["name"] == variant_name

        # Get
        response = api_client_with_organization.get(
            f"/items/{item_id}/variants/{variant_id}"
        )
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["id"] == variant_id

        # Update
        new_name = f"{variant_name}_updated"
        response = api_client_with_organization.put(
            f"/items/{item_id}/variants/{variant_id}",
            {"name": new_name},
        )
        assert response.status_code == 200, response.text
        data = response.json()["data"]
        assert data["id"] == variant_id
        assert data["name"] == new_name

        # Delete
        response = api_client_with_organization.delete(
            f"/items/{item_id}/variants/{variant_id}"
        )
        assert response.status_code == 204


class TestCells:
    @pytest.fixture
    def cell_group(
        self,
        api_client_with_organization: APIClient,
        organization_unit: dict,
    ) -> dict:
        cell_group_name = str(uuid.uuid4())
        alias = generate_random_string()
        response = api_client_with_organization.post(
            "/cells-groups",
            {
                "name": cell_group_name,
                "alias": alias,
                "unitId": organization_unit["id"],
            },
        )
        assert response.status_code == 200, response.text
        return response.json()["data"]

    def test_full_cell_group_lifecycle(
        self,
        api_client_with_organization: APIClient,
        organization_unit: dict,
    ) -> None:
        # Create
        cell_group_name = str(uuid.uuid4())
        alias = generate_random_string()
        response = api_client_with_organization.post(
            "/cells-groups",
            {
                "name": cell_group_name,
                "alias": alias,
                "unitId": organization_unit["id"],
            },
        )
        assert response.status_code == 200, response.text
        cell_group_data = response.json()["data"]
        cell_group_id = cell_group_data["id"]
        assert cell_group_data["name"] == cell_group_name
        assert cell_group_data["alias"] == alias
        assert cell_group_data["unitId"] == organization_unit["id"]

        # Get
        response = api_client_with_organization.get(f"/cells-groups/{cell_group_id}")
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["id"] == cell_group_id
        assert data["name"] == cell_group_name
        assert data["alias"] == alias
        assert data["unitId"] == organization_unit["id"]

        # List
        response = api_client_with_organization.get("/cells-groups")
        assert response.status_code == 200
        data = response.json()["data"]
        filtered = (x for x in data if x["id"] == cell_group_id)
        cell_group = next(filtered)
        assert cell_group is not None
        assert cell_group["name"] == cell_group_name
        assert cell_group["alias"] == alias
        assert cell_group["unitId"] == organization_unit["id"]

        # Update
        new_name = f"{cell_group_name}_updated"
        new_alias = f"{alias}D"
        response = api_client_with_organization.put(
            f"/cells-groups/{cell_group_id}",
            {"name": new_name, "alias": new_alias, "unitId": organization_unit["id"]},
        )
        assert response.status_code == 200, response.text

        # Delete
        response = api_client_with_organization.delete(f"/cells-groups/{cell_group_id}")
        assert response.status_code == 204

    def test_full_cell_lifecycle(
        self,
        api_client_with_organization: APIClient,
        cell_group: dict,
    ) -> None:
        # Create
        alias = generate_random_string()
        response = api_client_with_organization.post(
            f"/cells-groups/{cell_group['id']}/cells",
            {"alias": alias, "row": 1, "level": 1, "position": 1},
        )
        assert response.status_code == 200, response.text
        data = response.json()["data"]
        cell_id = data["id"]
        assert data["alias"] == alias
        assert data["row"] == 1
        assert data["level"] == 1
        assert data["position"] == 1

        # List
        response = api_client_with_organization.get(
            f"/cells-groups/{cell_group['id']}/cells"
        )
        assert response.status_code == 200
        data = response.json()["data"]
        assert len(data) == 1
        assert data[0]["id"] == cell_id

        # Get
        response = api_client_with_organization.get(f"/cells/{cell_id}")
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["id"] == cell_id
        assert data["alias"] == alias
        assert data["row"] == 1
        assert data["level"] == 1
        assert data["position"] == 1

        # Update
        new_alias = f"{alias}D"
        response = api_client_with_organization.put(
            f"/cells/{cell_id}",
            {"alias": new_alias, "row": 1, "level": 1, "position": 1},
        )
        assert response.status_code == 200, response.text
        data = response.json()["data"]
        assert data["id"] == cell_id
        assert data["alias"] == new_alias
        assert data["row"] == 1
        assert data["level"] == 1
        assert data["position"] == 1

        # Delete
        response = api_client_with_organization.delete(f"/cells/{cell_id}")
        assert response.status_code == 204


class TestInstance:
    def test_full_instance_lifecycle(
        self,
        api_client_with_organization: APIClient,
        organization_unit: dict,
    ) -> None:
        # Create Cells group
        cell_group_name = str(uuid.uuid4())
        alias = generate_random_string()
        response = api_client_with_organization.post(
            "/cells-groups",
            {
                "name": cell_group_name,
                "alias": alias,
                "unitId": organization_unit["id"],
            },
        )
        assert response.status_code == 200, response.text
        cell_group_data = response.json()["data"]

        # Create Cells
        alias_1 = generate_random_string()
        response = api_client_with_organization.post(
            f"/cells-groups/{cell_group_data['id']}/cells",
            {"alias": alias_1, "row": 1, "level": 1, "position": 1},
        )
        assert response.status_code == 200, response.text
        cell_data_1 = response.json()["data"]

        alias_2 = generate_random_string()
        response = api_client_with_organization.post(
            f"/cells-groups/{cell_group_data['id']}/cells",
            {"alias": alias_2, "row": 2, "level": 1, "position": 1},
        )
        assert response.status_code == 200, response.text
        cell_data_2 = response.json()["data"]

        # Create Item
        item_name = str(uuid.uuid4())
        response = api_client_with_organization.post(
            "/items",
            {"name": item_name},
        )
        assert response.status_code == 200, response.text
        item_data = response.json()["data"]

        # Create Variant
        variant_name = str(uuid.uuid4())
        article = str(random.randint(100000000000, 999999999999))
        response = api_client_with_organization.post(
            f"/items/{item_data['id']}/variants",
            {"name": variant_name, "article": article},
        )
        assert response.status_code == 200, response.text
        variant_data = response.json()["data"]

        # Create Instance
        response = api_client_with_organization.post(
            f"/items/{item_data['id']}/instances",
            data={
                "variantId": variant_data["id"],
                "cellId": cell_data_1["id"],
            },
        )
        assert response.status_code == 200, response.text
        instance_data = response.json()["data"]
        assert instance_data["status"] == "available"
        assert instance_data["variant"]["id"] == variant_data["id"]
        assert instance_data["variant"]["name"] == variant_name
        assert instance_data["variant"]["article"] == article

        assert instance_data["cell"]["id"] == cell_data_1["id"]
        assert instance_data["cell"]["alias"] == alias_1
        assert instance_data["cell"]["row"] == 1
        assert instance_data["cell"]["level"] == 1
        assert instance_data["cell"]["position"] == 1

        # create task
        response = api_client_with_organization.post(
            "/tasks",
            data={
                "name": "Перемещение",
                "description": "Перемещение в ячейку 2",
                "instanceId": instance_data["id"],
                "type": "movement",
                "unitId": organization_unit["id"],
                "items": [
                    {
                        "instanceId": instance_data["id"],
                        "targetCellId": cell_data_2["id"],
                    }
                ],
            },
        )

        assert response.status_code == 200, print(response.text)
