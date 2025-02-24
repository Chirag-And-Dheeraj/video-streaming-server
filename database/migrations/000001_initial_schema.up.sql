CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS videos (
    video_id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    upload_initiate_time TIMESTAMP,
    upload_status SMALLINT CHECK (upload_status IN (0, 1)),
    upload_end_time TIMESTAMP,
    user_id TEXT,
    delete_flag SMALLINT CHECK (delete_flag IN (0, 1)),
    FOREIGN KEY (user_id) REFERENCES users(id)
);