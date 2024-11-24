CREATE TABLE IF NOT EXISTS tags
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS novels_tags
(
    novel_id INTEGER NOT NULL,
    tag_id   INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS heroes
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS hero_moods
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR NOT NULL,
    image   VARCHAR NOT NULL,
    hero_id INTEGER NOT NULL,
    FOREIGN KEY (hero_id) REFERENCES heroes (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sound_transitions
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    sound    VARCHAR NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS backgrounds
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    image    VARCHAR NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS background_sounds
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    sound    VARCHAR NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);
