BEGIN;
    CREATE TYPE role_type_new AS ENUM (
        'normal',
        'moderator',
        'admin',
        'patreon',
    );
    UPDATE roles SET role_type = 'normal' WHERE role_type = 'sar';
    ALTER TABLE roles 
        ALTER COLUMN roles TYPE role_type_new;
            USING (roles::text::role_type_new)
    DROP TYPE role_type;
    ALTER TYPE role_type_new RENAME TO role_type;
COMMIT;
