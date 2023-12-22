-- +goose Up

CREATE TABLE IF NOT EXISTS guilds
(
    guildid      VARCHAR(25) NOT NULL DEFAULT '' PRIMARY KEY,
    autovoiceids TEXT        NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS permissions
(
    roleid  VARCHAR(25) NOT NULL DEFAULT '' PRIMARY KEY,
    guildid VARCHAR(25) NOT NULL DEFAULT '',
    perms   TEXT        NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS votes
(
    id       VARCHAR(25) NOT NULL DEFAULT '' PRIMARY KEY,
    jsondata TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS autovoice
(
    userid   VARCHAR(25) NOT NULL DEFAULT '' PRIMARY KEY,
    jsondata TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS refreshtokens
(
    ident   VARCHAR(25) NOT NULL DEFAULT '' PRIMARY KEY,
    token   TEXT        NOT NULL DEFAULT '',
    expires TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS apitokens
(
    ident       VARCHAR(25) NOT NULL DEFAULT '' PRIMARY KEY,
    salt        TEXT        NOT NULL DEFAULT '',
    created     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires     TIMESTAMP,
    lastaccess TIMESTAMP,
    hits        INT         NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS users
(
    userid       VARCHAR(25) NOT NULL DEFAULT '' PRIMARY KEY,
    redditoptout BOOLEAN     NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS reddit
(
    id      SERIAL PRIMARY KEY,
    guildid TEXT   NOT NULL DEFAULT '',
    userid  TEXT   NOT NULL DEFAULT '',
    karma   BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS redditsettings
(
    guildid   VARCHAR(25) PRIMARY KEY,
    state     BOOLEAN     NOT NULL DEFAULT TRUE,
    emotesinc TEXT        NOT NULL DEFAULT '',
    emotesdec TEXT        NOT NULL DEFAULT '',
    tokens    BIGINT      NOT NULL DEFAULT 1,
    penalty   INTEGER     NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS redditblocklist
(
    id      SERIAL PRIMARY KEY,
    userid  VARCHAR(25) NOT NULL DEFAULT '',
    guildid VARCHAR(25) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS redditrules
(
    id       VARCHAR(25) PRIMARY KEY,
    guildid  VARCHAR(25) NOT NULL DEFAULT '',
    trigger  INTEGER     NOT NULL DEFAULT 0,
    value    INTEGER     NOT NULL DEFAULT 0,
    action   VARCHAR(30) NOT NULL DEFAULT '',
    argument TEXT        NOT NULL DEFAULT '',
    checksum TEXT        NOT NULL DEFAULT ''
);

-- +goose Down

DROP TABLE IF EXISTS guilds;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS autovoice;
DROP TABLE IF EXISTS refreshtokens;
DROP TABLE IF EXISTS apitokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS reddit;
DROP TABLE IF EXISTS redditsettings;
DROP TABLE IF EXISTS redditblocklist;
DROP TABLE IF EXISTS redditrules;