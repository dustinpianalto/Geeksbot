BEGIN;
    ALTER TABLE users
        ALTER COLUMN active DROP NOT NULL,
        ALTER COLUMN staff DROP NOT NULL,
        ALTER COLUMN admin DROP NOT NULL;
    ALTER TABLE channels
        ALTER COLUMN guild_id DROP NOT NULL,
        ALTER COLUMN admin DROP NOT NULL,
        ALTER COLUMN default_channel DROP NOT NULL,
        ALTER COLUMN new_patron DROP NOT NULL;
    ALTER TABLE messages
        ALTER COLUMN created_at DROP NOT NULL,
        ALTER COLUMN content DROP NOT NULL,
        ALTER COLUMN channel_id DROP NOT NULL,
        ALTER COLUMN author_id DROP NOT NULL;
    ALTER TABLE messages
        ADD COLUMN embed json,
        ADD COLUMN previous_embeds json[];
    ALTER TABLE patreon_creators
        ALTER COLUMN creator DROP NOT NULL,
        ALTER COLUMN link DROP NOT NULL,
        ALTER COLUMN guild_id DROP NOT NULL;
    ALTER TABLE patreon_creators
        RENAME TO patreon_creator;
    ALTER TABLE patreon_tiers
        ALTER COLUMN name DROP NOT NULL,
        ALTER COLUMN creator DROP NOT NULL,
        ALTER COLUMN role DROP NOT NULL;
    ALTER TABLE patreon_tiers
        RENAME TO patreon_tier;
    ALTER TABLE requests
        ALTER COLUMN author_id DROP NOT NULL,
        ALTER COLUMN channel_id DROP NOT NULL,
        ALTER COLUMN content DROP NOT NULL,
        ALTER COLUMN requested_at DROP NOT NULL,
        ALTER COLUMN completed DROP NOT NULL,
        ALTER COLUMN message_id DROP NOT NULL,
        ALTER COLUMN guild_id DROP NOT NULL;
    ALTER TABLE roles
        ALTER COLUMN role_type DROP NOT NULL,
        ALTER COLUMN guild_id DROP NOT NULL;
    ALTER TABLE servers
        ALTER COLUMN name DROP NOT NULL,
        ALTER COLUMN ip_address DROP NOT NULL,
        ALTER COLUMN port DROP NOT NULL,
        ALTER COLUMN password DROP NOT NULL,
        ALTER COLUMN guild_id DROP NOT NULL;
COMMIT;
