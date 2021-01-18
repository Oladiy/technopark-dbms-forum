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
    created TIMESTAMP WITH TIME ZONE DEFAULT Now(),
    path INTEGER [] DEFAULT '{0}':: INTEGER []
);

CREATE TABLE Vote (
    id SERIAL PRIMARY KEY,
    nickname citext NOT NULL,
    voice INTEGER NOT NULL
);

CREATE OR REPLACE FUNCTION update_path()
    RETURNS TRIGGER AS
$BODY$
DECLARE
    parent_path         INT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := ARRAY(SELECT NEW.id::INTEGER);
    ELSE
        SELECT thread, path
        FROM Post
        WHERE thread = NEW.thread AND id = NEW.parent
        INTO first_parent_thread, parent_path;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Parent not exists' USING ERRCODE = '00404';
        END IF ;
        NEW.path := parent_path || NEW.id;
    END IF;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER path_updater
    BEFORE INSERT
    ON Post
    FOR EACH ROW
EXECUTE PROCEDURE update_path();
