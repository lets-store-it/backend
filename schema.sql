CREATE TYPE task_type AS ENUM ('movement', 'pickment');
CREATE TYPE task_status AS ENUM ('pending', 'in_progress', 'ready', 'completed', 'cancelled');
CREATE TYPE task_item_status AS ENUM ('pending', 'picked', 'done', 'returned', 'canceled');
CREATE TYPE item_instance_status AS ENUM ('available', 'reserved', 'consumed');


CREATE TABLE app_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,

    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255),

    yandex_id VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX app_user_email_idx ON app_user(email);

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

    width INTEGER CHECK (width > 0), -- in mm
    depth INTEGER CHECK (depth > 0), -- in mm
    height INTEGER CHECK (height > 0), -- in mm
    weight INTEGER CHECK (weight > 0), -- in g

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
    ean13 BIGINT CHECK (ean13::text ~ '^[0-9]{13}$'),

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

    type task_type NOT NULL,
    status task_status NOT NULL DEFAULT 'pending',
    
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    
    assigned_to_user_id UUID REFERENCES app_user(id),
    assigned_at TIMESTAMP,
    completed_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX task_org_id_idx ON task(org_id);
CREATE INDEX task_status_idx ON task(status) WHERE deleted_at IS NULL;
CREATE INDEX task_assigned_user_idx ON task(assigned_to_user_id, status) WHERE deleted_at IS NULL;
CREATE INDEX task_unit_id_idx ON task(unit_id);

CREATE TABLE item_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    item_id UUID NOT NULL REFERENCES item(id) ON DELETE RESTRICT,
    variant_id UUID NOT NULL REFERENCES item_variant(id) ON DELETE RESTRICT,

    cell_id UUID REFERENCES cell(id) ON DELETE RESTRICT,
    status item_instance_status NOT NULL DEFAULT 'available',
    affected_by_task_id UUID REFERENCES task(id) ON DELETE SET NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX item_instance_status_idx ON item_instance(status) WHERE deleted_at IS NULL;
CREATE INDEX item_instance_task_idx ON item_instance(affected_by_task_id) WHERE affected_by_task_id IS NOT NULL;
CREATE INDEX item_instance_cell_idx ON item_instance(cell_id) WHERE cell_id IS NOT NULL;

CREATE TABLE task_item (
    org_id UUID NOT NULL REFERENCES org(id),
    task_id UUID NOT NULL REFERENCES task(id) ON DELETE CASCADE,
    item_instance_id UUID NOT NULL REFERENCES item_instance(id) ON DELETE RESTRICT,
    status task_item_status NOT NULL DEFAULT 'pending',
    source_cell_id UUID REFERENCES cell(id) ON DELETE RESTRICT,
    destination_cell_id UUID REFERENCES cell(id) ON DELETE RESTRICT,
    PRIMARY KEY (task_id, item_instance_id)
);
CREATE INDEX task_item_instance_idx ON task_item(item_instance_id);
CREATE INDEX task_item_source_cell_idx ON task_item(source_cell_id) WHERE source_cell_id IS NOT NULL;
CREATE INDEX task_item_dest_cell_idx ON task_item(destination_cell_id) WHERE destination_cell_id IS NOT NULL;


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
    (8, 'rbac', 'employee'),
    (9, 'tasks', 'task'),
    (10, 'items', 'variant'),
    (11, 'org', 'api-token');


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
CREATE INDEX app_object_change_org_id_idx ON app_object_change(org_id);
CREATE INDEX app_object_change_target_object_type_idx ON app_object_change(target_object_type);
CREATE INDEX app_object_change_target_object_id_idx ON app_object_change(target_object_id);

CREATE TABLE app_api_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    name VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP
);
CREATE INDEX app_api_token_org_id_idx ON app_api_token(org_id);

CREATE TABLE tv_board (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org(id),
    unit_id UUID NOT NULL REFERENCES org_unit(id),
    name VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX tv_board_org_id_idx ON tv_board(org_id);
CREATE INDEX tv_board_token_idx ON tv_board(token);
