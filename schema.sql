CREATE TABLE org (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
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

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX item_org_id_idx ON item(org_id, id);

CREATE TABLE item_variant (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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

CREATE TABLE item_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES item(id),
    variant_id UUID NOT NULL REFERENCES item_variant(id),

    -- cell_id UUID REFERENCES cell(id),
    -- status VARCHAR(255) NOT NULL CHECK (status IN ('available', 'reserved', 'consumed')),
    -- affected_by_operation_id UUID REFERENCES operation(id)

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (item_id, variant_id)
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

CREATE TABLE app_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,

    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255),

    yandex_id VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE object_type (
--     id INTEGER PRIMARY KEY,
--     group VARCHAR(100) NOT NULL,
--     name VARCHAR(100) NOT NULL,
--     UNIQUE (group, name)
-- );
-- CREATE INDEX object_type_id_idx ON object_type(id);

-- INSERT INTO object_type (id, group, name) VALUES 
--     (1, 'storage', 'group'),
--     (2, 'storage', 'cells-group'),
--     (3, 'storage', 'cell'),
--     (4, 'items', 'item'),
--     (5, 'items', 'instance');

-- CREATE TABLE app_role (
--     id INTEGER PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     display_name VARCHAR(255) NOT NULL,
--     description VARCHAR(255)
-- );

-- INSERT INTO app_role (id, name, display_name) VALUES
--     (1, 'org_owner', 'Владелец', 'Имеет полный доступ к Организации, может назначать Управляющих организацией'),
--     (2, 'org_admin', 'Управляющий', 'Имеет полный доступ к Организации, может назначать Менеджеров организации'),
--     (3, 'org_manager', 'Менеджер', 'Имеет доступ к управлению объектами Организации, может назначать Сотрудников склада'),
--     (4, 'org_worker', 'Сотрудник', 'Имеет доступ к складким операциям'),

-- CREATE TABLE app_role_binding (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     role_id UUID NOT NULL REFERENCES app_role(id),
--     user_id UUID NOT NULL REFERENCES app_user(id),
--     UNIQUE (role_id, employee_id)
-- );
 
-- CREATE TABLE app_role_permission (
--     id INTEGER PRIMARY KEY,
--     permission VARCHAR(255) NOT NULL,
--     UNIQUE (role_id, permission)
-- );

-- CREATE TABLE role_permission (
--     role_id UUID NOT NULL REFERENCES role(id),
--     permission VARCHAR(255) NOT NULL,
--     UNIQUE (role_id, permission)
-- );

-- CREATE TABLE object_type (
--     id INTEGER PRIMARY KEY,
--     group VARCHAR(100) NOT NULL,
--     name VARCHAR(100) NOT NULL,
--     UNIQUE (group, name)
-- );
-- CREATE INDEX object_type_id_idx ON object_type(id);

-- INSERT INTO object_type (id, group, name) VALUES 
--     (1, 'storage', 'group'),
--     (2, 'storage', 'cells-group'),
--     (3, 'storage', 'cell'),
--     (4, 'items', 'item'),
--     (5, 'items', 'instance');

-- CREATE TABLE custom_field (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     org_id UUID NOT NULL REFERENCES org(id),
--     type VARCHAR(100) NOT NULL CHECK (type IN ('text', 'integer', 'decimal' 'boolean')),
--     name VARCHAR(100) NOT NULL,
--     label VARCHAR(100) NOT NULL,
--     description VARCHAR(255),
--     group_name VARCHAR(100),
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP,
--     UNIQUE (org_id, name)
-- );
-- CREATE INDEX custom_field_org_id_idx ON custom_field(org_id, id);
-- CREATE INDEX custom_field_group_name_idx ON custom_field(org_id, name);
-- -- relate

-- CREATE TABLE custom_field_related_types (
--     custom_field_id UUID NOT NULL REFERENCES custom_field(id),
--     object_type_id INTEGER NOT NULL REFERENCES object_type(id),
--     PRIMARY KEY (custom_field_id, object_type_id)
-- );
-- CREATE INDEX custom_field_related_types_custom_field_id_idx ON custom_field_related_types(custom_field_id);
-- CREATE INDEX custom_field_related_types_object_type_id_idx ON custom_field_related_types(object_type_id);


-- CREATE TABLE cell_kind (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     org_id UUID NOT NULL REFERENCES org(id),
--     name VARCHAR(255) NOT NULL,
--     height INTEGER NOT NULL,
--     width INTEGER NOT NULL,
--     depth INTEGER NOT NULL,
--     max_weight INTEGER NOT NULL,
--     UNIQUE (org_id, name)
-- );

-- CREATE TABLE cell_group (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     space_id UUID NOT NULL REFERENCES storage_group(id),
--     name VARCHAR(255) NOT NULL,
--     short_name VARCHAR(255),
--     UNIQUE (space_id, name)
-- );

-- CREATE TABLE cell (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     short_name VARCHAR(255) NOT NULL,
--     group_id UUID NOT NULL REFERENCES cell_group(id),
--     kind_id UUID NOT NULL REFERENCES cell_kind(id),
--     rack VARCHAR(255) NOT NULL,
--     level INTEGER NOT NULL,
--     position INTEGER NOT NULL,
--     UNIQUE (group_id, rack, level, position)
-- );

-- CREATE TABLE item (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     org_id UUID NOT NULL REFERENCES org(id),
--     name VARCHAR(255) NOT NULL,
--     ean VARCHAR(255) NOT NULL,
--     UNIQUE (org_id, ean)
-- );

-- CREATE TABLE item_property (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     item_id UUID NOT NULL REFERENCES item(id),
--     name VARCHAR(255) NOT NULL,
--     type VARCHAR(255) NOT NULL CHECK (type IN ('text', 'bool', 'number')),
--     value VARCHAR(255) NOT NULL,
--     UNIQUE (item_id, name)
-- );

-- CREATE TABLE item_variant (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     item_id UUID NOT NULL REFERENCES item(id),
--     name VARCHAR(255) NOT NULL,
--     width INTEGER,
--     depth INTEGER,
--     height INTEGER,
--     weight INTEGER,
--     UNIQUE (item_id, name),
--     CHECK ((width IS NULL AND depth IS NULL AND height IS NULL AND weight IS NULL) OR 
--            (width IS NOT NULL AND depth IS NOT NULL AND height IS NOT NULL AND weight IS NOT NULL))
-- );

-- CREATE TABLE item_variant_property (
--     variant_id UUID NOT NULL REFERENCES item_variant(id),
--     property_id UUID NOT NULL REFERENCES item_property(id),
--     value VARCHAR(255) NOT NULL,
--     UNIQUE (variant_id, property_id)
-- );

-- CREATE TABLE item_instance (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     tracking_id VARCHAR(255) NOT NULL UNIQUE,
--     item_id UUID NOT NULL REFERENCES item(id),
--     variant_id UUID NOT NULL REFERENCES item_variant(id),
--     cell_id UUID REFERENCES cell(id),
--     status VARCHAR(255) NOT NULL CHECK (status IN ('available', 'reserved', 'consumed')),
--     affected_by_operation_id UUID REFERENCES operation(id)
-- );

-- CREATE TABLE app_user (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     fullname VARCHAR(255),
--     email VARCHAR(255) NOT NULL UNIQUE
-- );

-- CREATE TABLE employee (
--     user_id UUID NOT NULL REFERENCES app_user(id),
--     org_id UUID NOT NULL REFERENCES org(id),
--     PRIMARY KEY (user_id, org_id)
-- );

-- CREATE TABLE role (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     name VARCHAR(255) NOT NULL,
--     display_name VARCHAR(255) NOT NULL
-- );

-- CREATE TABLE role_permission (
--     role_id UUID NOT NULL REFERENCES role(id),
--     permission VARCHAR(255) NOT NULL,
--     UNIQUE (role_id, permission)
-- );

-- CREATE TABLE role_binding (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     role_id UUID NOT NULL REFERENCES role(id),
--     employee_id UUID NOT NULL REFERENCES employee(id),
--     UNIQUE (role_id, employee_id)
-- );

-- CREATE TABLE operation (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     type VARCHAR(255) NOT NULL CHECK (type IN ('pick', 'movement')),
--     assigned_to UUID REFERENCES employee(id),
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );

-- CREATE TABLE operation_item (
--     operation_id UUID NOT NULL REFERENCES operation(id),
--     item_instance_id UUID NOT NULL REFERENCES item_instance(id),
--     status VARCHAR(255) NOT NULL CHECK (status IN ('pending', 'picked', 'done', 'returned')),
--     origin_cell_id UUID REFERENCES cell(id),
--     destination_cell_id UUID REFERENCES cell(id)
-- );

