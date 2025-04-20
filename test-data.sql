INSERT INTO org (id, name, subdomain) VALUES
    ('453f0e17-c8f4-4c99-9d20-f0e13572550e', 'ООО "TapMe"', 'tapme');

INSERT INTO org_unit (id, org_id, name, alias, address) VALUES
    ('ba548db7-ceae-4ff6-bc30-653f45a53ddd', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'Москва Амурская', 'AM', 'ул. Амурская, 1'),
    ('c3ac04a4-043d-4b34-b2bd-4b7e138ca3ad', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'Красноярск 1', 'KR1', NULL);

INSERT INTO storage_group (id, org_id, unit_id, parent_id, name, alias) VALUES
    ('73717d25-ba60-459a-bf8e-736ba9fa3e60', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'ba548db7-ceae-4ff6-bc30-653f45a53ddd', NULL, 'Основная группа', 'AOG'),
    ('3fa3f6a0-5c9b-4360-bf2d-110a7acd6745', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'c3ac04a4-043d-4b34-b2bd-4b7e138ca3ad', NULL, 'Первый Этаж', 'KFS'),
    ('3726e287-b6a5-47b3-b178-11f99f860fe3', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'c3ac04a4-043d-4b34-b2bd-4b7e138ca3ad', '3fa3f6a0-5c9b-4360-bf2d-110a7acd6745', 'Зона А', 'KZA');

INSERT INTO cells_group (id, org_id, storage_group_id, name, alias) VALUES
    ('c58851e0-0233-4d20-b965-48992ef8f9fc', '453f0e17-c8f4-4c99-9d20-f0e13572550e', '3726e287-b6a5-47b3-b178-11f99f860fe3', 'Стеллаж 1', 'SFJ1');

INSERT INTO cell (id, org_id, cells_group_id, alias, row, level, position) VALUES
    ('65f7de59-0848-48b2-b541-1a072d7ae461', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'c58851e0-0233-4d20-b965-48992ef8f9fc', 'A11', 1, 1, 1),
    ('8ee7ee39-8471-4537-98d8-0da9ae25043e', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'c58851e0-0233-4d20-b965-48992ef8f9fc', 'A12', 1, 2, 1),
    ('77dea1ec-6ce9-45f5-abae-3dc27fc23b82', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'c58851e0-0233-4d20-b965-48992ef8f9fc', 'A13', 1, 3, 1);

INSERT INTO item (id, org_id, name, description, width, depth, height, weight) VALUES
    ('5b29eb11-c699-414b-a351-32476ad9dae0', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'Визитка TapMe Карта', 'Визитка в виде карты', 100, 100, 100, 100),
    ('5985325b-2e81-4d4a-8081-b96f034c81bb', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'Визитка TapMe Стикер', 'Визитка в виде стикера', 100, 100, 100, 100);

INSERT INTO item_variant (id, item_id, org_id, name, article, ean13) VALUES
    ('29ad8e0e-f187-47c7-a158-01de64a2aad3', '5985325b-2e81-4d4a-8081-b96f034c81bb', '453f0e17-c8f4-4c99-9d20-f0e13572550e','Черная', 'TM01', NULL),
    ('7513a498-9e6f-48fc-9787-e6395df65a85', '5985325b-2e81-4d4a-8081-b96f034c81bb', '453f0e17-c8f4-4c99-9d20-f0e13572550e','Белая', 'TM02', NULL),

    ('99921684-433c-4ebe-a461-d5a9d19519b3', '5b29eb11-c699-414b-a351-32476ad9dae0', '453f0e17-c8f4-4c99-9d20-f0e13572550e','Основной', NULL, NULL)

INSERT INTO app_user (id, email, first_name, last_name, middle_name) VALUES
    ('71254334-e6d0-47bb-8483-d13cc8456910', 'evseev-2003@yandex.ru', 'Евгений', 'Евсеев', 'Васильевич');

INSERT INTO app_role_binding (id, org_id, role_id, user_id) VALUES
    ('094ae0c3-23fb-4c0d-bd36-283f040eb8a3', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 1, '71254334-e6d0-47bb-8483-d13cc8456910');

INSERT INTO app_api_token (id, org_id, name, token) VALUES
    ('b826bb73-6018-4342-a122-66f418442ba2', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'testin-token', 'super-secret-token');
