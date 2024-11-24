BEGIN; 

CREATE TABLE IF NOT EXISTS groups (
    group_id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS songs (
    song_id bigserial PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL UNIQUE,
    lyrics TEXT NOT NULL,
    release_date DATE,
    link VARCHAR(255),
    group_id BIGINT NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups (group_id) ON DELETE CASCADE
);


CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);
CREATE INDEX IF NOT EXISTS idx_songs_group_id ON songs(group_id);
CREATE INDEX IF NOT EXISTS idx_songs_song ON songs(song);
CREATE INDEX IF NOT EXISTS idx_songs_release_date ON songs(release_date);


COMMIT;