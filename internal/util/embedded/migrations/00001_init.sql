-- +goose Up

CREATE TABLE IF NOT EXISTS guilds
(
    guild_id      VARCHAR(25) NOT NULL DEFAULT '',
    autorole_ids  text        NOT NULL DEFAULT '',
    autovoice_ids text        NOT NULL DEFAULT '',
    PRIMARY KEY (guild_id)
);

CREATE TABLE IF NOT EXISTS permissions
(
    role_id  VARCHAR(25) NOT NULL DEFAULT '',
    guild_id VARCHAR(25) NOT NULL DEFAULT '',
    perms    text        NOT NULL DEFAULT '',
    PRIMARY KEY (role_id)
);

CREATE TABLE IF NOT EXISTS votes
(
    id        VARCHAR(25) NOT NULL DEFAULT '',
    json_data TEXT        NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS autovoice
(
    user_id        VARCHAR(25) NOT NULL DEFAULT '',
    json_data TEXT        NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS guildapi
(
    guild_id      VARCHAR(25) NOT NULL DEFAULT '',
    enabled       BOOLEAN     NOT NULL DEFAULT TRUE,
    origins       TEXT        NOT NULL DEFAULT '',
    token_hash    TEXT        NOT NULL DEFAULT '',
    PRIMARY KEY (guild_id)
);

CREATE TABLE IF NOT EXISTS refreshtokens
(
    user_id       VARCHAR(25) NOT NULL DEFAULT '',
    refresh_token TEXT        NOT NULL DEFAULT '',
    expires_at    TIMESTAMP      NOT NULL DEFAULT '1970-01-01 00:00:00',
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS apitokens
(
    user_id             VARCHAR(25) NOT NULL DEFAULT '',
    salt                TEXT        NOT NULL DEFAULT '',
    created_at          TIMESTAMP      NOT NULL DEFAULT '1970-01-01 00:00:00',
    expires_at          TIMESTAMP      NOT NULL DEFAULT '1970-01-01 00:00:00',
    last_accessed_at    TIMESTAMP      NOT NULL DEFAULT '1970-01-01 00:00:00',
    hits                INT      NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id)
);

-- +goose Down

DROP TABLE IF EXISTS guilds;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS autovoice;
DROP TABLE IF EXISTS guildapi;
DROP TABLE IF EXISTS refreshtokens;