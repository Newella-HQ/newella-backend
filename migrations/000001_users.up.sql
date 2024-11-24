CREATE TABLE IF NOT EXISTS users
(
    id          SERIAL PRIMARY KEY,
    sub_id      VARCHAR(255) UNIQUE NOT NULL, -- from JWT
    username    VARCHAR UNIQUE      NOT NULL, -- default = email
    real_name   VARCHAR             NOT NULL, -- jwt
    description VARCHAR,
    email       VARCHAR(255) UNIQUE,
    picture     VARCHAR,
    role        VARCHAR(30)         NOT NULL DEFAULT 'user',
    created_at  TIMESTAMP                    DEFAULT now()
);

CREATE TABLE IF NOT EXISTS oauth_tokens
(
    access_token  VARCHAR NOT NULL, -- oauth
    refresh_token VARCHAR NOT NULL, -- oauth
    user_id       INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users_subs
(
    user_id       INTEGER NOT NULL,
    subscriber_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (subscriber_id) REFERENCES users (id) ON DELETE CASCADE
);
