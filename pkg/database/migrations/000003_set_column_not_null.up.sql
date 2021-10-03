BEGIN;
    ALTER TABLE users
        ALTER COLUMN active SET NOT NULL,
        ALTER COLUMN staff SET NOT NULL,
        ALTER COLUMN admin SET NOT NULL;
    ALTER TABLE channels
        ALTER COLUMN guild_id SET NOT NULL,
        ALTER COLUMN admin SET NOT NULL,
        ALTER COLUMN default_channel SET NOT NULL,
        ALTER COLUMN new_patron SET NOT NULL;
    ALTER TABLE messages
        ALTER COLUMN created_at SET NOT NULL,
        ALTER COLUMN content SET NOT NULL,
        ALTER COLUMN channel_id SET NOT NULL,
        ALTER COLUMN author_id SET NOT NULL;
    ALTER TABLE messages
        DROP COLUMN embed,
        DROP COLUMN previous_embeds;
    ALTER TABLE patreon_creator
        ALTER COLUMN creator SET NOT NULL,
        ALTER COLUMN link SET NOT NULL,
        ALTER COLUMN guild_id SET NOT NULL;
    ALTER TABLE patreon_creator
        RENAME TO patreon_creators;
    ALTER TABLE patreon_tier
        ALTER COLUMN name SET NOT NULL,
        ALTER COLUMN creator SET NOT NULL,
        ALTER COLUMN role SET NOT NULL;
    ALTER TABLE patreon_tier
        RENAME TO patreon_tiers;
    ALTER TABLE requests
        ALTER COLUMN author_id SET NOT NULL,
        ALTER COLUMN channel_id SET NOT NULL,
        ALTER COLUMN content SET NOT NULL,
        ALTER COLUMN requested_at SET NOT NULL,
        ALTER COLUMN completed SET NOT NULL,
        ALTER COLUMN message_id SET NOT NULL,
        ALTER COLUMN guild_id SET NOT NULL;
    ALTER TABLE roles
        ALTER COLUMN role_type SET NOT NULL,
        ALTER COLUMN guild_id SET NOT NULL;
    ALTER TABLE servers
        ALTER COLUMN name SET NOT NULL,
        ALTER COLUMN ip_address SET NOT NULL,
        ALTER COLUMN port SET NOT NULL,
        ALTER COLUMN password SET NOT NULL,
        ALTER COLUMN guild_id SET NOT NULL;
COMMIT;
