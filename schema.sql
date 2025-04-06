CREATE TABLE org (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    subdomain VARCHAR(255) NOT NULL UNIQUE CHECK (subdomain ~ '^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$'),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX org_subdomain_idx ON org(subdomain);

CREATE TABLE org_unit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    UNIQUE (org_id, name)
);

CREATE TABLE storage_space (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unit_id UUID NOT NULL REFERENCES org_unit(id),
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(255),
    UNIQUE (unit_id, name)
);

CREATE TABLE cell_kind (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    name VARCHAR(255) NOT NULL,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    depth INTEGER NOT NULL,
    max_weight INTEGER NOT NULL,
    UNIQUE (org_id, name)
);

CREATE TABLE cell_group (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id UUID NOT NULL REFERENCES storage_space(id),
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(255),
    UNIQUE (space_id, name)
);

CREATE TABLE cell (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_name VARCHAR(255) NOT NULL,
    group_id UUID NOT NULL REFERENCES cell_group(id),
    kind_id UUID NOT NULL REFERENCES cell_kind(id),
    rack VARCHAR(255) NOT NULL,
    level INTEGER NOT NULL,
    position INTEGER NOT NULL,
    UNIQUE (group_id, rack, level, position)
);

CREATE TABLE item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    name VARCHAR(255) NOT NULL,
    ean VARCHAR(255) NOT NULL,
    UNIQUE (org_id, ean)
);

CREATE TABLE item_property (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES item(id),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL CHECK (type IN ('text', 'bool', 'number')),
    value VARCHAR(255) NOT NULL,
    UNIQUE (item_id, name)
);

CREATE TABLE item_variant (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES item(id),
    name VARCHAR(255) NOT NULL,
    width INTEGER,
    depth INTEGER,
    height INTEGER,
    weight INTEGER,
    UNIQUE (item_id, name),
    CHECK ((width IS NULL AND depth IS NULL AND height IS NULL AND weight IS NULL) OR 
           (width IS NOT NULL AND depth IS NOT NULL AND height IS NOT NULL AND weight IS NOT NULL))
);

CREATE TABLE item_variant_property (
    variant_id UUID NOT NULL REFERENCES item_variant(id),
    property_id UUID NOT NULL REFERENCES item_property(id),
    value VARCHAR(255) NOT NULL,
    UNIQUE (variant_id, property_id)
);

CREATE TABLE item_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tracking_id VARCHAR(255) NOT NULL UNIQUE,
    item_id UUID NOT NULL REFERENCES item(id),
    variant_id UUID NOT NULL REFERENCES item_variant(id),
    cell_id UUID REFERENCES cell(id),
    status VARCHAR(255) NOT NULL CHECK (status IN ('available', 'reserved', 'consumed')),
    affected_by_operation_id UUID REFERENCES operation(id)
);

CREATE TABLE app_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE employee (
    user_id UUID NOT NULL REFERENCES app_user(id),
    org_id UUID NOT NULL REFERENCES org(id),
    PRIMARY KEY (user_id, org_id)
);

CREATE TABLE role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL
);

CREATE TABLE role_permission (
    role_id UUID NOT NULL REFERENCES role(id),
    permission VARCHAR(255) NOT NULL,
    UNIQUE (role_id, permission)
);

CREATE TABLE role_binding (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES role(id),
    employee_id UUID NOT NULL REFERENCES employee(id),
    UNIQUE (role_id, employee_id)
);

CREATE TABLE operation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(255) NOT NULL CHECK (type IN ('pick', 'movement')),
    assigned_to UUID REFERENCES employee(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE operation_item (
    operation_id UUID NOT NULL REFERENCES operation(id),
    item_instance_id UUID NOT NULL REFERENCES item_instance(id),
    status VARCHAR(255) NOT NULL CHECK (status IN ('pending', 'picked', 'done', 'returned')),
    origin_cell_id UUID REFERENCES cell(id),
    destination_cell_id UUID REFERENCES cell(id)
);

