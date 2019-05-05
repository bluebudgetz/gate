/*
 * This file is called by "gate" to add initial data to the database.
 *
 * It IS ONLY invoked in development environment.
 */
INSERT INTO bb.accounts (id, created_on, updated_on, deleted_on, name, parent_id)
VALUES (1, DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), NULL, NULL, 'Bank1', NULL),
       (2, DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), NULL, NULL, 'Bank2', NULL),
       (3, DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), NULL, NULL, 'Loans', 1),
       (4, DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), NULL, NULL, 'Salary', 1),
       (5, DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), NULL, NULL, 'Insurance', 2),
       (6, DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), NULL, NULL, 'CheapAndHealthy Inc.', 5),
       (7, DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), NULL, NULL, 'SafeStruct Ltd.', 5);

INSERT INTO bb.transactions (created_on, origin, source_account_id, target_account_id, amount, comments)
VALUES (DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), 'Init', 4, 6, 100, 'Health insurance'),
       (DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), 'Init', 4, 7, 35, 'Structure insurance'),
       (DATE_ADD(DATE(NOW()), INTERVAL FLOOR(RAND() * 24) HOUR), 'Init', 4, 3, 1269, 'Pay loan 1');
