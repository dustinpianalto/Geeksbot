BEGIN;
    ALTER TABLE messages
        ALTER COLUMN previous_content DROP NOT NULL;
    ALTER TABLE messages
        ALTER COLUMN previous_content DROP DEFAULT;
COMMIT;
