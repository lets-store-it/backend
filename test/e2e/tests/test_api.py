import random
import string
import uuid

import pytest

from tests.apiclient.apiclient import APIClient

API_BASE = "http://localhost:8080"
SESSION_COOKIE = "6964d3a3-8bb3-406e-9e0f-1324f3aa868f"


@pytest.fixture(scope="session")
def api_client() -> APIClient:
    client = APIClient(API_BASE, SESSION_COOKIE)
    response = client.get("/me")
    assert response.status_code == 200
    return client


@pytest.fixture(scope="session")
def api_client_with_organization(api_client: APIClient) -> APIClient:
    org_name = str(uuid.uuid4())
    response = api_client.post("/orgs", {"name": org_name, "subdomain": org_name})
    assert response.status_code == 200
    org_data = response.json()["data"]

    client = APIClient(API_BASE, SESSION_COOKIE)
    response = client.get("/me")
    assert response.status_code == 200
    client.set_organization(org_data["id"])
    return client


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


def generate_random_string(length: int = 4) -> str:
    return "".join(random.choices(string.ascii_letters, k=length))


class TestOrganization:
    def test_full_organization_lifecycle(self, api_client: APIClient) -> None:
        # Create
        org_name = str(uuid.uuid4())
        response = api_client.post("/orgs", {"name": org_name, "subdomain": org_name})
        assert response.status_code == 200
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
        assert response.status_code == 200
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
        new_name = f"{new_name}_patched"
        response = api_client_with_organization.patch(
            f"/units/{unit_id}",
            {"name": new_name},
        )
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["alias"] == new_alias
        assert data["address"] == new_address

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

        # Patch
        new_name = f"{new_name}_patched"
        response = api_client_with_organization.patch(
            f"/storage-groups/{storage_group_id}",
            {"name": new_name},
        )
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["alias"] == new_alias
        assert data["unitId"] == organization_unit["id"]

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
        print(data)
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

        # Patch
        new_name = f"{new_name}_patched"
        response = api_client_with_organization.patch(
            f"/cells-groups/{cells_group_id}",
            {"name": new_name},
        )
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["alias"] == new_alias
        assert data["unitId"] == organization_unit["id"]

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
        print(item_data)
        assert item_data["name"] == name
        assert item_data["description"] == description

        # List
        response = api_client_with_organization.get("/items")
        assert response.status_code == 200
        data = response.json()["data"]
        print(data)
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

        # Patch
        new_name = f"{new_name}_patched"
        response = api_client_with_organization.patch(
            f"/items/{item_id}",
            {"name": new_name},
        )
        assert response.status_code == 200
        data = response.json()["data"]
        assert data["name"] == new_name
        assert data["description"] == new_description

        # Delete
        response = api_client_with_organization.delete(f"/items/{item_id}")
        assert response.status_code == 204
