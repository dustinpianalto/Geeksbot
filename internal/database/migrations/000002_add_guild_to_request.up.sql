ALTER TABLE requests 
    ADD COLUMN guild_id varchar(30) REFERENCES guilds(id) ON DELETE CASCADSE;
