CREATE TABLE IF NOT EXISTS genres
(
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    parent_id INTEGER,
    FOREIGN KEY (parent_id) REFERENCES genres (id) ON DELETE CASCADE
);

CREATE TYPE novel_length AS ENUM ('small', 'normal', 'large');

CREATE TABLE IF NOT EXISTS novels
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR      NOT NULL,
    description VARCHAR,
    price       INTEGER      NOT NULL DEFAULT 0,
    length      novel_length NOT NULL,
    image       VARCHAR, -- may be null
    rating      FLOAT        NOT NULL,
    genre_id    INTEGER      NOT NULL,
    user_id     VARCHAR(255) NOT NULL,
    FOREIGN KEY (genre_id) REFERENCES genres (id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);
