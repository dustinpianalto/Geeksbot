BEGIN;
    CREATE TABLE IF NOT EXISTS guilds (
        id varchar(30),
        new_patron_message varchar(1000),
        prefixes varchar(10)[],
        PRIMARY KEY(id)
    );
    CREATE TYPE role_type as ENUM (
        'normal',
        'moderator',
        'admin',
        'patreon'
    );
    CREATE TABLE IF NOT EXISTS roles (
        id varchar(30),
        role_type role_type,
        guild_id varchar(30),
        PRIMARY KEY(id),
        CONSTRAINT fk_guild
            FOREIGN KEY(guild_id)
                REFERENCES guilds(id)
                ON DELETE CASCADE
    );
    CREATE TABLE IF NOT EXISTS channels (
        id varchar(30),
        guild_id varchar(30),
        admin boolean,
        default_channel boolean,
        new_patron boolean,
        PRIMARY KEY(id),
        CONSTRAINT fk_guild
            FOREIGN KEY(guild_id)
                REFERENCES guilds(id)
                ON DELETE CASCADE
    );
    CREATE TABLE IF NOT EXISTS users (
        id varchar(30),
        steam_id varchar(30),
        active boolean,
        staff boolean,
        admin boolean,
        PRIMARY KEY(id)
    );
    CREATE TABLE IF NOT EXISTS messages (
        id varchar(30),
        created_at timestamp,
        modified_at timestamp,
        content varchar(2000),
        previous_content varchar(2000)[],
        channel_id varchar(30),
        author_id varchar(30),
        embed json,
        previous_embeds json[],
        PRIMARY KEY(id),
        CONSTRAINT fk_channel
            FOREIGN KEY(channel_id)
                REFERENCES channels(id)
                ON DELETE CASCADE,
        CONSTRAINT fk_user
            FOREIGN KEY(author_id)
                REFERENCES users(id)
                ON DELETE SET NULL
    );
    CREATE TABLE IF NOT EXISTS patreon_creator (
        id integer GENERATED ALWAYS AS IDENTITY,
        creator varchar(100),
        link varchar(200),
        guild_id varchar(30),
        PRIMARY KEY(id),
        CONSTRAINT fk_guild
            FOREIGN KEY(guild_id)
                REFERENCES guilds(id)
                ON DELETE CASCADE
    );
    CREATE TABLE IF NOT EXISTS patreon_tier (
        id  integer GENERATED ALWAYS AS IDENTITY,
        name varchar(100),
        description varchar(1000),
        creator integer,
        role varchar(30),
        next_tier integer,
        PRIMARY KEY(id),
        CONSTRAINT fk_creator
            FOREIGN KEY(creator)
                REFERENCES patreon_creator(id)
                ON DELETE CASCADE,
        CONSTRAINT fk_role
            FOREIGN KEY(role)
                REFERENCES roles(id)
                ON DELETE SET NULL,
        CONSTRAINT fk_tier
            FOREIGN KEY(next_tier)
                REFERENCES patreon_tier(id)
                ON DELETE SET NULL
    );
    CREATE TABLE IF NOT EXISTS requests (
        id integer GENERATED ALWAYS AS IDENTITY,
        author_id varchar(30),
        channel_id varchar(30),
        content varchar(1000),
        requested_at timestamp,
        completed_at timestamp,
        completed boolean,
        completed_by varchar(30),
        message_id varchar(30),
        completed_message varchar(1000),
        PRIMARY KEY(id),
        CONSTRAINT fk_user
            FOREIGN KEY(author_id)
                REFERENCES users(id)
                ON DELETE CASCADE,
        CONSTRAINT fk_channel
            FOREIGN KEY(channel_id)
                REFERENCES channels(id)
                ON DELETE CASCADE,
        CONSTRAINT fk_completed_by
            FOREIGN KEY(completed_by)
                REFERENCES users(id)
                ON DELETE SET NULL,
        CONSTRAINT fk_message
            FOREIGN KEY(message_id)
                REFERENCES messages(id)
                ON DELETE CASCADE
    );
    CREATE TABLE IF NOT EXISTS comments (
        id integer GENERATED ALWAYS AS IDENTITY,
        author_id varchar(30),
        request_id integer,
        comment_at timestamp,
        content varchar(1000),
        PRIMARY KEY(id),
        CONSTRAINT fk_user
            FOREIGN KEY(author_id)
                REFERENCES users(id)
                ON DELETE CASCADE,
        CONSTRAINT fk_request
            FOREIGN KEY(request_id)
                REFERENCES requests(id)
                ON DELETE CASCADE
    );
    CREATE TABLE IF NOT EXISTS servers (
        id integer GENERATED ALWAYS AS IDENTITY,
        name varchar(100),
        ip_address varchar(15),
        port integer,
        password varchar(200),
        alerts_channel_id varchar(30),
        guild_id varchar(30),
        info_channel_id varchar(30),
        info_message_id varchar(30),
        settings_message_id varchar(30),
        PRIMARY KEY(id),
        CONSTRAINT fk_alert_channel
            FOREIGN KEY(alerts_channel_id)
                REFERENCES channels(id)
                ON DELETE SET NULL,
        CONSTRAINT fk_guild
            FOREIGN KEY(guild_id)
                REFERENCES guilds(id)
                ON DELETE CASCADE,
        CONSTRAINT fk_info_channel
            FOREIGN KEY(info_channel_id)
                REFERENCES channels(id)
                ON DELETE SET NULL,
        CONSTRAINT fk_info_message
            FOREIGN KEY(info_message_id)
                REFERENCES messages(id)
                ON DELETE SET NULL,
        CONSTRAINT fk_settings_message
            FOREIGN KEY(settings_message_id)
                REFERENCES messages(id)
                ON DELETE SET NULL
    );
COMMIT;
