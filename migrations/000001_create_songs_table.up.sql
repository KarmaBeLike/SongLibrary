CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_title VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    release_date DATE,
    link VARCHAR(255)
);
