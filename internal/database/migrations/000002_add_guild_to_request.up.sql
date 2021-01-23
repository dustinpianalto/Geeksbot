ALTER TABLE requests 
    ADD COLUMN guild_id varchar(30) CONSTRAINT fk_guild REFERENCES guilds(id) ON DELETE CASCADE;
