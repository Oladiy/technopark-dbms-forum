CREATE EXTENSION IF NOT EXISTS citext;

SET timezone = 'Europe/Moscow';

CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    nickname citext NOT NULL UNIQUE,
    fullname VARCHAR(256) NOT NULL,
    about TEXT,
    email citext NOT NULL UNIQUE
);

CREATE TABLE Forum (
    id SERIAL PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    author citext NOT NULL,
    slug citext NOT NULL UNIQUE,
    posts INTEGER DEFAULT 0,
    threads INTEGER DEFAULT 0
);

CREATE TABLE Thread (
    id SERIAL PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    author citext NOT NULL,
    forum citext NOT NULL,
    message TEXT,
    votes INTEGER DEFAULT 0,
    slug citext NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT Now()
);

CREATE TABLE Post (
    id SERIAL PRIMARY KEY,
    parent INTEGER DEFAULT 0,
    author citext NOT NULL,
    message TEXT,
    isEdited BOOLEAN DEFAULT FALSE,
    forum citext NOT NULL,
    thread INTEGER,
    created TIMESTAMP WITH TIME ZONE DEFAULT Now()
);

CREATE TABLE Vote (
    id SERIAL PRIMARY KEY,
    nickname citext REFERENCES Users(nickname),
    voice INTEGER NOT NULL
);