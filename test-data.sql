INSERT INTO org (id, name, subdomain) VALUES
    ('a4da775e-dbd4-4e1c-a254-232a4686dcde', 'ООО "Экзотика"', 'exotic-inc'),
    ('453f0e17-c8f4-4c99-9d20-f0e13572550e', 'ООО "TapMe"', 'tapme');

INSERT INTO org_unit (id, org_id, name, alias) VALUES
    ('ba548db7-ceae-4ff6-bc30-653f45a53ddd', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'Москва Амурская', 'AM'),
    ('c3ac04a4-043d-4b34-b2bd-4b7e138ca3ad', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'Красноярск 1', 'KR1'),

    ('3b6e1ac8-66c2-4315-94f5-f6276eaefed9', 'a4da775e-dbd4-4e1c-a254-232a4686dcde', 'ГимН3', 'MOF');

brbr    ('73717d25-ba60-459a-bf8e-736ba9fa3e60', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'ba548db7-ceae-4ff6-bc30-653f45a53ddd', NULL, 'Входная зона', 'EN1'),
    ('3fa3f6a0-5c9b-4360-bf2d-110a7acd6745', 'a4da775e-dbd4-4e1c-a254-232a4686dcde', '3b6e1ac8-66c2-4315-94f5-f6276eaefed9', NULL, 'Основная группа', 'PS'),
    ('3726e287-b6a5-47b3-b178-11f99f860fe3', 'a4da775e-dbd4-4e1c-a254-232a4686dcde', '3b6e1ac8-66c2-4315-94f5-f6276eaefed9', '3fa3f6a0-5c9b-4360-bf2d-110a7acd6745', 'Зона А', 'ZA');

INSERT INTO item (id, org_id, name, description) VALUES
    ('5b29eb11-c699-414b-a351-32476ad9dae0', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'Визитка TapMe Карта', 'Визитка в виде карты'),
    ('5985325b-2e81-4d4a-8081-b96f034c81bb', '453f0e17-c8f4-4c99-9d20-f0e13572550e', 'Визитка TapMe Стикер', 'Визитка в виде стикера'),

    ('97e4d2cf-8519-4cf9-aff3-ca8368506fd2', 'a4da775e-dbd4-4e1c-a254-232a4686dcde', 'Давитель для корма', 'Давитель для корма');

INSERT INTO item_variant (id, item_id, name, article, ean13) VALUES
    ('29ad8e0e-f187-47c7-a158-01de64a2aad3', '97e4d2cf-8519-4cf9-aff3-ca8368506fd2', 'Черная', 'TM01', NULL),
    ('7513a498-9e6f-48fc-9787-e6395df65a85', '97e4d2cf-8519-4cf9-aff3-ca8368506fd2', 'Белая', 'TM02', NULL),

    ('99921684-433c-4ebe-a461-d5a9d19519b3', '97e4d2cf-8519-4cf9-aff3-ca8368506fd2', 'Основной', NULL, NULL)