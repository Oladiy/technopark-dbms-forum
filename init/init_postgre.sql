ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET max_worker_processes = '4';
ALTER SYSTEM SET max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET max_parallel_workers = '4';
ALTER SYSTEM SET max_parallel_maintenance_workers = '2';
ALTER SYSTEM SET random_page_cost = '0.1';
ALTER SYSTEM SET seq_page_cost = '0.1';
ALTER SYSTEM SET wal_buffers = '6912kB';

CREATE EXTENSION IF NOT EXISTS citext;

SET timezone = 'Europe/Moscow';

CREATE UNLOGGED TABLE Users (
    id SERIAL PRIMARY KEY,
    nickname citext NOT NULL UNIQUE,
    fullname VARCHAR(256) NOT NULL,
    about TEXT,
    email citext NOT NULL UNIQUE
);

CREATE INDEX index_users_full_info ON Users(nickname, fullname, about, email);

CREATE UNLOGGED TABLE Forum (
    id SERIAL PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    author citext NOT NULL,
    slug citext NOT NULL UNIQUE,
    posts INTEGER DEFAULT 0,
    threads INTEGER DEFAULT 0
);

CREATE INDEX index_forum_author ON Forum(author);
CREATE INDEX index_forum_full_info ON Forum(title, author, slug, posts, threads);

CREATE UNLOGGED TABLE ForumUsers (
    id SERIAL PRIMARY KEY,
    forum_slug citext NOT NULL REFERENCES Forum(slug),
    user_nickname citext NOT NULL REFERENCES Users(nickname),
    UNIQUE(forum_slug, user_nickname)
);

CREATE UNLOGGED TABLE Thread (
    id SERIAL PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    author citext NOT NULL,
    forum citext NOT NULL,
    message TEXT,
    votes INTEGER DEFAULT 0,
    slug citext NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT Now()
);

CREATE INDEX index_thread_author ON Thread(author);
CREATE INDEX index_thread_forum ON Thread(Forum);
CREATE INDEX index_thread_slug ON Thread(slug);
CREATE INDEX index_thread_full_info ON Thread(title, author, forum, message, slug, votes, created);

CREATE UNLOGGED TABLE Post (
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

CREATE UNLOGGED TABLE Vote (
    id SERIAL PRIMARY KEY,
    nickname citext NOT NULL,
    voice INTEGER NOT NULL,
    thread INTEGER NOT NULL,
    UNIQUE(nickname, thread)
);

CREATE INDEX index_vote_nickname_voice ON Vote(nickname, voice);

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
BEGIN
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
