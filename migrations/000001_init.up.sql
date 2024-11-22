CREATE TABLE IF NOT EXISTS "user"
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR     NOT NULL,
    real_name     VARCHAR     NOT NULL,
    description   VARCHAR     NOT NULL,
    email         VARCHAR(255) UNIQUE,
    password_hash VARCHAR     NOT NULL,
    is_confirmed  BOOLEAN     NOT NULL DEFAULT FALSE,
    role          VARCHAR(30) NOT NULL DEFAULT 'user',
    created_at    TIMESTAMP            DEFAULT now()
);

CREATE TABLE IF NOT EXISTS novel
(
    id          SERIAL  NOT NULL,
    name        VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    price       INTEGER NOT NULL,
    length      INTEGER NOT NULL,
    status      VARCHAR NOT NULL,
    image       VARCHAR NOT NULL,
    rating      FLOAT   NOT NULL,
    genre_id    INTEGER NOT NULL,
    author_id   INTEGER NOT NULL,
    FOREIGN KEY (genre_id) REFERENCES genre (id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS genre
(
    id        SERIAL PRIMARY KEY,
    name      VARCHAR NOT NULL,
    parent_id VARCHAR NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES genre (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_subs
(
    user_id       INTEGER NOT NULL,
    subscriber_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (subscriber_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS library
(
    id       SERIAL PRIMARY KEY,
    type     VARCHAR NOT NULL,
    user_id  INTEGER NOT NULL,
    novel_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS progress
(
    id                 SERIAL PRIMARY KEY,
    slide_group_number INTEGER NOT NULL,
    slide_number       INTEGER NOT NULL,
    library_id         INTEGER NOT NULL,
    FOREIGN KEY (library_id) REFERENCES library (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rating
(
    id       SERIAL PRIMARY KEY,
    rating   VARCHAR NOT NULL,
    user_id  INTEGER NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novel (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment
(
    id       SERIAL PRIMARY KEY,
    text     VARCHAR NOT NULL,
    user_id  INTEGER NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novel (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS complaint
(
    id          SERIAL PRIMARY KEY,
    subject     VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    status      VARCHAR NOT NULL,
    user_id     INTEGER NOT NULL,
    novel_id    INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novel (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tag
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS novel_tag
(
    novel_id INTEGER NOT NULL,
    tag_id   INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novel (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tag (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS hero
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novel (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS hero_mood
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR NOT NULL,
    image   VARCHAR NOT NULL,
    hero_id INTEGER NOT NULL,
    FOREIGN KEY (hero_id) REFERENCES hero (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS sound_transition
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    sound    VARCHAR NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novel (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS background
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    image    VARCHAR NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novel (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS background_sound
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR NOT NULL,
    sound    VARCHAR NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES novel (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS slide
(
    id                  SERIAL PRIMARY KEY,
    name                VARCHAR NOT NULL,
    replica             VARCHAR NOT NULL,
    background_id       INTEGER NOT NULL,
    hero_mood_id        INTEGER NOT NULL,
    sound_transition_id INTEGER NOT NULL,
    slide_group_id      INTEGER NOT NULL,
    background_sound_id INTEGER NOT NULL,
    FOREIGN KEY (background_id) REFERENCES background (id) ON DELETE CASCADE,
    FOREIGN KEY (hero_mood_id) REFERENCES hero_mood (id) ON DELETE CASCADE,
    FOREIGN KEY (sound_transition_id) REFERENCES sound_transition (id) ON DELETE CASCADE,
    FOREIGN KEY (slide_group_id) REFERENCES slide_group (id) ON DELETE CASCADE,
    FOREIGN KEY (background_sound_id) REFERENCES background_sound (id) ON DELETE CASCADE

);
CREATE TABLE IF NOT EXISTS slide_group
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR NOT NULL,
    novel_id       INTEGER NOT NULL,
    start_slide_id INTEGER NOT NULL,
    end_slide_id   INTEGER NOT NULL,
    FOREIGN KEY (start_slide_id) REFERENCES slide (id) ON DELETE CASCADE,
    FOREIGN KEY (end_slide_id) REFERENCES slide (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS answer
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR NOT NULL,
    link_to    INTEGER NOT NULL,
    belongs_id INTEGER NOT NULL,
    FOREIGN KEY (link_to) REFERENCES slide (id) ON DELETE CASCADE,
    FOREIGN KEY (belongs_id) REFERENCES slide (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS custom_filter
(
    id        SERIAL PRIMARY KEY,
    type      VARCHAR NOT NULL,
    min_value INTEGER NOT NULL,
    max_value INTEGER NOT NULL,
    sort_type INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS filter_value
(
    id               SERIAL PRIMARY KEY,
    value            VARCHAR NOT NULL,
    sequence_number  INTEGER NOT NULL,
    custom_filter_id INTEGER NOT NULL,
    FOREIGN KEY (custom_filter_id) REFERENCES custom_filter (id) ON DELETE CASCADE
);