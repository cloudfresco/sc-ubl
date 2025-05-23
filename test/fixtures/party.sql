INSERT INTO `parties` (uuid4,
    party_name,
    registration_name,
    registration_date,
    registration_expiration_date,
    company_id,
    level_p,
    parent_id,
    num_chd,
    leaf,
    address_id,
    status_code,
    created_at,
    updated_at,
    created_by_user_id,
    updated_by_user_id) VALUES (UNHEX(REPLACE('ccabf4e9-c992-4c91-b008-e4a26138dd1c','-','')), 'Consortial','','2023-03-24 10:04:26','2024-07-25 10:04:26','',0,0,1,true,1,'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');

INSERT INTO `parties` (uuid4,
    party_name,
    registration_name,
    registration_date,
    registration_expiration_date,
    company_id,
    level_p,
    parent_id,
    num_chd,
    leaf,
    address_id,
    status_code,
    created_at,
    updated_at,
    created_by_user_id,
    updated_by_user_id) VALUES (UNHEX(REPLACE('bba18409-f347-44ff-8a2a-f3def98be3cc','-','')),'IYT Corporation','','2023-03-24 10:04:26','2024-07-25 10:04:26','',0,0,0,false,2,'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');
   
INSERT INTO `parties` (uuid4,
    party_name,
    registration_name,
    registration_date,
    registration_expiration_date,
    company_id,
    level_p,
    parent_id,
    num_chd,
    leaf,
    address_id,
    status_code,
    created_at,
    updated_at,
    created_by_user_id,
    updated_by_user_id) VALUES (UNHEX(REPLACE('32144a7e-f37a-41af-8df6-cebe40d88252','-','')),'Salescompany ltd','The Sellercompany Incorporated','2023-03-24 10:04:26','2024-07-25 10:04:26','5402697509',0,0,0,false,3,'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');
 
 INSERT INTO `parties` (uuid4,
    party_name,
    registration_name,
    registration_date,
    registration_expiration_date,
    company_id,
    level_p,
    parent_id,
    num_chd,
    leaf,
    address_id,
    status_code,
    created_at,
    updated_at,
    created_by_user_id,
    updated_by_user_id) VALUES (UNHEX(REPLACE('f71ef0a6-d878-4f2f-af7e-e424fe6ef4d0','-','')),'Buyercompany ltd','The buyercompany inc.','2023-03-24 10:04:26','2024-07-25 10:04:26','5645342123',0,0,0,false,4,'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');

INSERT INTO `parties` (uuid4,
    party_name,
    registration_name,
    registration_date,
    registration_expiration_date,
    company_id,
    level_p,
    parent_id,
    num_chd,
    leaf,
    address_id,
    status_code,
    created_at,
    updated_at,
    created_by_user_id,
    updated_by_user_id) VALUES (UNHEX(REPLACE('f720686f-85d3-435a-8865-06e03a5818a9','-','')),'Ebeneser Scrooge Inc','Ebeneser Scrooge Inc.','2023-03-24 10:04:26','2024-07-25 10:04:26','6411982340',0,0,0,false,5,'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');

INSERT INTO `parties` (uuid4,
    party_name,
    registration_name,
    registration_date,
    registration_expiration_date,
    company_id,
    level_p,
    parent_id,
    num_chd,
    leaf,
    address_id,
    status_code,
    created_at,
    updated_at,
    created_by_user_id,
    updated_by_user_id) VALUES (UNHEX(REPLACE('976e86b5-a959-4717-a00a-25a717f155c3','-','')),'Test supplier','','2023-03-24 10:04:26','2024-07-25 10:04:26','',0,0,0,false,6,'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');

INSERT INTO `parties` (uuid4,
    party_name,
    registration_name,
    registration_date,
    registration_expiration_date,
    company_id,
    level_p,
    parent_id,
    num_chd,
    leaf,
    address_id,
    status_code,
    created_at,
    updated_at,
    created_by_user_id,
    updated_by_user_id) VALUES (UNHEX(REPLACE('ee0c5f7b-1c3e-4287-b3af-07b4c9707f9b','-','')),'Test customer','','2023-03-24 10:04:26','2024-07-25 10:04:26','',1,1,0,false,7,'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');

INSERT INTO `party_chds` (uuid4,
   party_id,
   party_chd_id,
   status_code,
   created_at,
   updated_at,
   created_by_user_id,
   updated_by_user_id) VALUES(UNHEX(REPLACE('64fbb4ef-7268-45c6-98e4-d14eebea062e','-','')),1, 7,'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');
