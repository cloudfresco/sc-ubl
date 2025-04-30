INSERT INTO `despatch_lines`
	  (uuid4,
    note,
    line_status_code,
    delivered_quantity,
    backorder_quantity,
    item_id,
despatch_header_id,
status_code,
created_at,
updated_at,
created_by_user_id,
updated_by_user_id) VALUES (UNHEX(REPLACE('bd96e929-f6e8-46ab-932e-638e8ba52421','-','')),'Mrs Green agreed to waive charge','NoStatus',90,10, 8, 1, 'active','2019-07-23 10:04:26','2019-07-23 10:04:26','auth0|673c75d516e8adb9e6ffc892','auth0|673c75d516e8adb9e6ffc892');
