/*
 * This file is called by "docker-compose.yml" to create the database.
 *
 * It IS NOT used by the actual servers.
 */
CREATE OR REPLACE DATABASE bb DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
GRANT ALL ON bb.* TO 'bb'@'%' IDENTIFIED BY 'bb';
