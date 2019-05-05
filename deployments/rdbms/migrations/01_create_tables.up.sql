USE bb;

CREATE TABLE accounts
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_on TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_on TIMESTAMP    NULL,
    name       VARCHAR(255) NOT NULL,
    parent_id  INT UNSIGNED NULL,
    CONSTRAINT FOREIGN KEY fk_accounts_parent_id (parent_id) REFERENCES accounts (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);
CREATE TABLE transactions
(
    id                INT UNSIGNED   NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_on        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    origin            VARCHAR(255)   NOT NULL,
    source_account_id INT UNSIGNED   NOT NULL,
    target_account_id INT UNSIGNED   NOT NULL,
    amount            FLOAT UNSIGNED NOT NULL,
    comments          TEXT,
    CONSTRAINT FOREIGN KEY fk_transactions_source_account_id (source_account_id) REFERENCES accounts (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,
    CONSTRAINT FOREIGN KEY fk_transactions_target_account_id (target_account_id) REFERENCES accounts (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);
