INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_coord_lat,
    location_coord_lon,
    location_type_code,
    address_id
) VALUES (
    UNHEX(REPLACE('c703277f-84ca-4816-9ccf-fad8e202d3b6','-','')),
    'Hamburg',
    '53.551° N',
    '9.9937° E',
    'DEHAM',
    1
);
    
INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_type_code
) VALUES (
  UNHEX(REPLACE('84bfcf2e-403b-11eb-bc4a-1fc4aa7d879d','-','')),
  'The Factory',
  'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_type_code
) VALUES (
    UNHEX(REPLACE('286c605e-4043-11eb-9c0b-7b4196cf71fa','-','')),
    'Port of Singapore',
    'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_type_code
) VALUES (
    UNHEX(REPLACE('770b7624-403d-11eb-b44b-d3f4ad185386','-','')),
    'Port of Rotterdam',
    'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_type_code
) VALUES  (
    UNHEX(REPLACE('770b7624-403d-11eb-b44b-d3f4ad185387','-','')),
    'Genneb',
    'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_type_code
) VALUES (
    UNHEX(REPLACE('770b7624-403d-11eb-b44b-d3f4ad185388','-','')),
    'Nijmegen',
    'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_type_code
) VALUES (
    UNHEX(REPLACE('7f29ce3c-403d-11eb-9579-6bd2f4cf4ed6','-','')),
    'The Warehouse',
    'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    address_id,
    location_coord_lat,
    location_coord_lon,
    location_type_code
) VALUES (
    UNHEX(REPLACE('01670315-a51f-4a11-b947-ce8e245128eb','-','')),
    'Lagkagehuset Islands Brygge',
    2,
    '55.6642249',
    '12.57341045',
    'USNYC'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_coord_lat,
    location_coord_lon,
    location_type_code
) VALUES (
    UNHEX(REPLACE('b4454ae5-dcd4-4955-8080-1f986aa5c6c3','-','')),
    'Copenhagen',
    '55.671° N',
    '12.453° E',
    'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_coord_lat,
    location_coord_lon,
    location_type_code
) VALUES (
    UNHEX(REPLACE('1d09e9e9-dba3-4de1-8ef8-3ab6d32dbb40','-','')),
    'Orlando',
    '28.481° N',
    '-81.48° E',
    'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    location_coord_lat,
    location_coord_lon,
    location_type_code
) VALUES (
    UNHEX(REPLACE('ea9af21d-8471-47ac-aa59-e949ea74b08e','-','')),
    'Miami',
    '25.782° N',
    '-80.36° E',
    'USMIA'
);

INSERT INTO `locations` (
    uuid4,
    loc_name,
    address_id,
    location_coord_lat,
    location_coord_lon,
    location_type_code
) VALUES (
    UNHEX(REPLACE('06aca2f6-f1d0-48f8-ba46-9a3480adfd23','-','')),
    'Eiffel Tower',
    3,
    '48.8585500',
    '2.294492036',
    'USNYC'
);
