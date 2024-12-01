CREATE TYPE user_role AS ENUM ('user', 'moderator', 'admin');

CREATE TABLE IF NOT EXISTS users
(
    id          VARCHAR(255) PRIMARY KEY, -- sub_id from JWT Google
    username    VARCHAR UNIQUE NOT NULL,  -- default = email
    real_name   VARCHAR        NOT NULL,  -- jwt
    description VARCHAR,
    email       VARCHAR(255) UNIQUE,
    picture     VARCHAR,
    role        user_role      NOT NULL DEFAULT 'user',
    created_at  TIMESTAMP               DEFAULT now()
);

CREATE TABLE IF NOT EXISTS oauth_tokens
(
    access_token  VARCHAR      NOT NULL, -- oauth
    refresh_token VARCHAR      NOT NULL, -- oauth
    user_id       VARCHAR(255) NOT NULL UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users_subs
(
    user_id       VARCHAR(255) NOT NULL,
    subscriber_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (subscriber_id) REFERENCES users (id) ON DELETE CASCADE
);
