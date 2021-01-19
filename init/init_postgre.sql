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

CREATE TABLE ForumUsers (
    id SERIAL PRIMARY KEY,
    forum_slug citext NOT NULL REFERENCES Forum(slug),
    user_nickname citext NOT NULL REFERENCES Users(nickname),
    UNIQUE(forum_slug, user_nickname)
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

CREATE OR REPLACE FUNCTION update_forum_threads_counter()
RETURNS TRIGGER AS $update_forum_threads_counter$
BEGIN
    UPDATE Forum
    SET threads = threads + 1
    WHERE slug = new.forum;

    INSERT INTO ForumUsers(forum_slug, user_nickname)
    VALUES(new.forum, new.author)
    ON CONFLICT (forum_slug, user_nickname) DO NOTHING;

    RETURN NEW;
END;
$update_forum_threads_counter$ LANGUAGE plpgsql;

CREATE TRIGGER update_forum_threads_counter
    BEFORE INSERT ON Thread
    FOR EACH ROW EXECUTE PROCEDURE update_forum_threads_counter();

CREATE OR REPLACE FUNCTION update_forum_posts_counter()
RETURNS TRIGGER AS $update_forum_posts_counter$
begin
    UPDATE Forum
    SET posts = posts + 1
    WHERE slug = new.forum;

    INSERT INTO ForumUsers(forum_slug, user_nickname)
    VALUES (new.forum, new.author)
    ON CONFLICT (forum_slug, user_nickname) DO NOTHING;

    RETURN NEW;
END;
$update_forum_posts_counter$ LANGUAGE plpgsql;

CREATE TRIGGER update_forum_posts_counter
    BEFORE INSERT ON Post
    FOR EACH ROW EXECUTE PROCEDURE update_forum_posts_counter();
