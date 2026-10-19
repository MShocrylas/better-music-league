CREATE TABLE users (
    user_id  INTEGER PRIMARY KEY,
    username TEXT NOT NULL
) STRICT;

CREATE TABLE leagues (
    league_id INTEGER PRIMARY KEY,
    name      TEXT NOT NULL
) STRICT;

CREATE TABLE leagues_users (
    league_id INTEGER NOT NULL REFERENCES leagues(league_id) ON DELETE CASCADE,
    user_id   INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (league_id, user_id)
) STRICT;

CREATE TABLE rounds (
    round_id  INTEGER PRIMARY KEY,
    theme     TEXT NOT NULL,
    league_id INTEGER NOT NULL REFERENCES leagues(league_id) ON DELETE CASCADE
) STRICT;

CREATE TABLE songs (
    song_id INTEGER PRIMARY KEY,
    title   TEXT NOT NULL,
    artist  TEXT NOT NULL,
    url     TEXT,
    UNIQUE (title, artist)
) STRICT;

CREATE TABLE submissions (
    submission_id INTEGER PRIMARY KEY,
    user_id       INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, -- The user who submitted this song
    round_id      INTEGER NOT NULL REFERENCES rounds(round_id) ON DELETE CASCADE,
    song_id       INTEGER NOT NULL REFERENCES songs(song_id) ON DELETE RESTRICT,
    UNIQUE (user_id, round_id), -- A user can only have one submission per round
    UNIQUE (round_id, song_id) -- A song can only be submitted once per round
) STRICT;

CREATE TABLE votes (
    vote_id       INTEGER PRIMARY KEY,
    rank          INTEGER NOT NULL,
    submission_id INTEGER NOT NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    user_id       INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    UNIQUE (submission_id, user_id) -- User can only have one vote per submission
) STRICT;
