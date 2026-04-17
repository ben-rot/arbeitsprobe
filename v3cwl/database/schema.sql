-- database/schema.sql

CREATE TABLE IF NOT EXISTS users (

    id SERIAL PRIMARY KEY,

    discord_id VARCHAR(24) UNIQUE NOT NULL,
    username VARCHAR(32) NOT NULL,
    nickname VARCHAR(32) NOT NULL,
    avatar_hash VARCHAR(40) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW()
);



CREATE TABLE IF NOT EXISTS families (

    id SERIAL PRIMARY KEY,
    discord_server_id VARCHAR(24) NOT NULL UNIQUE,
    name VARCHAR(32) NOT NULL,
    server_icon VARCHAR(32) NOT NULL,

    registered_at TIMESTAMPTZ DEFAULT NOW()
);



CREATE TABLE IF NOT EXISTS clans (

    id SERIAL PRIMARY KEY,

    tag VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(32) NOT NULL,
    cwl_league INT NOT NULL,
    is_locked BOOLEAN DEFAULT FALSE,

    family_id INT REFERENCES families(id) ON DELETE CASCADE
);




CREATE TABLE IF NOT EXISTS accounts (

    id SERIAL PRIMARY KEY,

    tag VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(32) NOT NULL,
    townhall INT NOT NULL,

    owner_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    clan_id INT REFERENCES clans(id) ON DELETE SET NULL DEFAULT NULL,
    family_id INT REFERENCES families(id) ON DELETE SET NULL DEFAULT NULL
);




CREATE TABLE IF NOT EXISTS managers (

    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    family_id INT REFERENCES families(id) ON DELETE CASCADE,

    last_active_at TIMESTAMPTZ DEFAULT NULL,
    assigned_at TIMESTAMPTZ DEFAULT NOW()
);





CREATE INDEX IF NOT EXISTS idx_manager_family ON managers(family_id);
CREATE INDEX IF NOT EXISTS idx_account_family ON accounts(family_id);
CREATE INDEX IF NOT EXISTS idx_account_owner ON accounts(owner_id);
CREATE INDEX IF NOT EXISTS idx_account_clan ON accounts(clan_id);
CREATE INDEX IF NOT EXISTS idx_clan_family ON clans(family_id);
CREATE INDEX IF NOT EXISTS idx_families_covering_stats ON families(id) INCLUDE (discord_server_id, name, server_icon);