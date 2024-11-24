CREATE TABLE IF NOT EXISTS users_novels
(
    id       SERIAL PRIMARY KEY,
    type     VARCHAR NOT NULL, -- possible to replace varchar by enum
    user_id  INTEGER NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ratings
(
    id       SERIAL PRIMARY KEY,
    rating   FLOAT   NOT NULL,
    user_id  INTEGER NOT NULL,
    novel_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments
(
    id         SERIAL PRIMARY KEY,
    text       VARCHAR   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    user_id    INTEGER   NOT NULL,
    novel_id   INTEGER   NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS complaints
(
    id          SERIAL PRIMARY KEY,
    subject     VARCHAR     NOT NULL,
    description VARCHAR     NOT NULL,
    status      VARCHAR(50) NOT NULL,
    user_id     INTEGER     NOT NULL,
    novel_id    INTEGER     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);
