CREATE TABLE users (
    user_id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE COLLATE nocase
) STRICT;

CREATE TABLE leagues (
    league_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
) STRICT;

CREATE TABLE leagues_users (
    league_id INTEGER NOT NULL REFERENCES leagues(league_id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'member')),
    PRIMARY KEY (league_id, user_id)
) STRICT;
CREATE INDEX idx_leagues_users_user_id ON leagues_users (user_id);

CREATE TABLE rounds (
    round_id INTEGER PRIMARY KEY,
    league_id INTEGER NOT NULL REFERENCES leagues(league_id) ON DELETE CASCADE,
    theme TEXT NOT NULL,
    position INTEGER NOT NULL,
    submission_deadline TEXT CHECK (
        submission_deadline
        IS strftime('%Y-%m-%dT%H:%M:%SZ', submission_deadline)
    ),
    voting_deadline TEXT CHECK (
        voting_deadline IS strftime('%Y-%m-%dT%H:%M:%SZ', voting_deadline)
    ),
    -- Where in the order of rounds in a league does this round sit
    UNIQUE (league_id, position)
) STRICT;

CREATE TABLE songs (
    song_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    url TEXT,
    UNIQUE (title, artist)
) STRICT;

CREATE TABLE submissions (
    submission_id INTEGER PRIMARY KEY,
    song_id INTEGER NOT NULL REFERENCES songs(song_id) ON DELETE RESTRICT,
    -- The user who submitted this song
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    round_id INTEGER NOT NULL REFERENCES rounds(round_id) ON DELETE CASCADE,
    -- ISO-8601
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    UNIQUE (user_id, round_id), -- A user can only have one submission per round
    UNIQUE (round_id, song_id), -- A song can only be submitted once per round
    UNIQUE (submission_id, round_id) -- Exists as target for votes composite FK
) STRICT;

CREATE TABLE votes (
    vote_id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    submission_id INTEGER NOT NULL,
    round_id INTEGER NOT NULL,
    song_rank INTEGER NOT NULL CHECK (song_rank > 0),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    FOREIGN KEY (submission_id, round_id) REFERENCES submissions(
        submission_id, round_id
    ) ON DELETE CASCADE,
    -- User can only have one vote per submission
    UNIQUE (submission_id, user_id),
    -- Vote rank needs to be unique for this user for this round 
    -- (no repeated ranks on a ballot)
    UNIQUE (user_id, round_id, song_rank)
) STRICT;

CREATE TRIGGER votes_no_self_vote
BEFORE INSERT ON votes
WHEN
    new.user_id
    = (SELECT user_id FROM submissions WHERE submission_id = new.submission_id)
BEGIN
    SELECT raise(ABORT, 'cannot vote for own submission');
END;
