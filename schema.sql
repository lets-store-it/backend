CREATE TABLE app_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,

    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255),

    yandex_id VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE org (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(255) NOT NULL UNIQUE CHECK (subdomain ~ '^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$'),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX org_subdomain_idx ON org(subdomain);

CREATE TABLE org_unit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL,
    address VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (org_id, name),
    UNIQUE (org_id, alias)
);
CREATE INDEX org_unit_org_id_idx ON org_unit(org_id, id);
CREATE INDEX org_unit_alias_idx ON org_unit(org_id, alias);

CREATE TABLE storage_group (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    unit_id UUID NOT NULL REFERENCES org_unit(id),
    parent_id UUID,
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL,
    description VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (org_id, alias),
    FOREIGN KEY (parent_id) REFERENCES storage_group(id) ON DELETE CASCADE,
    CHECK (parent_id != id)
);
CREATE INDEX storage_group_org_id_idx ON storage_group(org_id, id);
CREATE INDEX storage_group_unit_id_idx ON storage_group(org_id, unit_id);

CREATE TABLE item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),

    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),

    width INTEGER, -- in mm
    depth INTEGER, -- in mm
    height INTEGER, -- in mm
    weight INTEGER, -- in g

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX item_org_id_idx ON item(org_id, id);

CREATE TABLE item_variant (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    item_id UUID NOT NULL REFERENCES item(id),

    name VARCHAR(255) NOT NULL,

    article VARCHAR(255),
    ean13 INTEGER,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (item_id, name)
);
CREATE INDEX item_variant_item_id_idx ON item_variant(item_id);
CREATE INDEX item_variant_ean13_idx ON item_variant(ean13);
CREATE INDEX item_variant_article_idx ON item_variant(article);


CREATE TABLE cells_group (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    unit_id UUID NOT NULL REFERENCES org_unit(id),
    storage_group_id UUID REFERENCES storage_group(id),
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    UNIQUE (org_id, alias)
);
CREATE INDEX cells_group_org_id_idx ON cells_group(org_id, id);

CREATE TABLE cell (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),

    cells_group_id UUID NOT NULL REFERENCES cells_group(id),
    
    alias VARCHAR(255) NOT NULL,

    row INTEGER NOT NULL,
    level INTEGER NOT NULL,
    position INTEGER NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    UNIQUE (cells_group_id, alias),
    UNIQUE (cells_group_id, row, level, position)
);
CREATE INDEX cell_cells_group_id_idx ON cell(cells_group_id, id);


CREATE TABLE task (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    unit_id UUID NOT NULL REFERENCES org_unit(id),

    type VARCHAR(255) NOT NULL CHECK (type IN ('pick', 'movement')),
    status VARCHAR(255) NOT NULL CHECK (status IN ('pending', 'in_progress', 'awaiting_to_collect', 'completed', 'failed')) DEFAULT 'pending',
    
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    
    assigned_to_user_id UUID REFERENCES app_user(id),
    assigned_at TIMESTAMP,
    completed_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE item_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    item_id UUID NOT NULL REFERENCES item(id),
    variant_id UUID NOT NULL REFERENCES item_variant(id),

    cell_id UUID REFERENCES cell(id),
    status VARCHAR(255) NOT NULL CHECK (status IN ('available', 'reserved', 'consumed')) DEFAULT 'available',
    affected_by_task_id UUID REFERENCES task(id),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (item_id, variant_id)
);

CREATE TABLE task_item (
    org_id UUID NOT NULL REFERENCES org(id),
    task_id UUID NOT NULL REFERENCES task(id),
    item_instance_id UUID NOT NULL REFERENCES item_instance(id),
    status VARCHAR(255) NOT NULL CHECK (status IN ('pending', 'picked', 'done', 'failed', 'returned')) DEFAULT 'pending',
    source_cell_id UUID REFERENCES cell(id),
    destination_cell_id UUID REFERENCES cell(id)
);


CREATE TABLE app_user_session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES app_user(id),
    token VARCHAR(255) NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP + INTERVAL '7 days',
    revoked_at TIMESTAMP
);
CREATE INDEX app_user_session_user_id_idx ON app_user_session(user_id);

CREATE TABLE app_role (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);

INSERT INTO app_role (id, name, display_name, description) VALUES
    (1, 'org_owner', 'Владелец', 'Имеет полный доступ к Организации, может назначать Управляющих организацией'),
    (2, 'org_admin', 'Управляющий', 'Имеет полный доступ к Организации, может назначать Менеджеров организации'),
    (3, 'org_manager', 'Менеджер', 'Имеет доступ к управлению объектами Организации, может назначать Сотрудников склада'),
    (4, 'org_worker', 'Сотрудник', 'Имеет доступ к складким операциям');

CREATE TABLE app_role_binding (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    role_id INTEGER NOT NULL REFERENCES app_role(id),
    user_id UUID NOT NULL REFERENCES app_user(id),
    UNIQUE (user_id, org_id)
);

CREATE TABLE object_type (
    id INTEGER PRIMARY KEY,
    object_group VARCHAR(100) NOT NULL,
    object_name VARCHAR(100) NOT NULL,
    UNIQUE (object_group, object_name)
);
CREATE INDEX object_type_id_idx ON object_type(id);

INSERT INTO object_type (id, object_group, object_name) VALUES 
    (1, 'org', 'organization'),
    (2, 'org', 'unit'),
    (3, 'storage', 'group'),
    (4, 'storage', 'cells-group'),
    (5, 'storage', 'cell'),
    (6, 'items', 'item'),
    (7, 'items', 'instance'),
    -- (8, 'rbac', 'user-roles'),
    (8, 'rbac', 'employee');


CREATE TABLE app_object_change (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    user_id UUID REFERENCES app_user(id),
    action VARCHAR(255) NOT NULL CHECK (action IN ('create', 'update', 'delete')),
    time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    target_object_type INTEGER NOT NULL REFERENCES object_type(id),
    target_object_id UUID NOT NULL,
    prechange_state JSONB,
    postchange_state JSONB
);

CREATE TABLE app_api_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    name VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP
);

CREATE TABLE tv_board (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    unit_id UUID NOT NULL REFERENCES org_unit(id),
    name VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
