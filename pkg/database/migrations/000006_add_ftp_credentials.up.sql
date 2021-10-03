BEGIN;
    ALTER TABLE servers
        ADD COLUMN ftp_port int4,
        ADD COLUMN ftp_username varchar(200),
        ADD COLUMN ftp_password varchar(200);
COMMIT;
