CREATE TABLE users (
    user_id BIGINT PRIMARY KEY UNIQUE
);

CREATE TABLE secrets (
    user_id BIGINT,
    FOREIGN KEY (user_id) REFERENCES users(user_id),

    secret_id SERIAL PRIMARY KEY,
    secret_name TEXT UNIQUE,
    secret_key TEXT UNIQUE
);
