BEGIN;
    ALTER TABLE servers
        DROP COLUMN ftp_port,
        DROP COLUMN ftp_username,
        DROP COLUMN ftp_password;
COMMIT;
