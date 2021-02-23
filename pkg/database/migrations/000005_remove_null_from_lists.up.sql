BEGIN;
    ALTER TABLE messages
        ALTER COLUMN previous_content SET NOT NULL;
    ALTER TABLE messages
        ALTER COLUMN prevous_content SET DEFAULT array[]::varchar[]
COMMIT;
