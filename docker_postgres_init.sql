CREATE TABLE IF NOT EXISTS artists(
      id          SERIAL   PRIMARY KEY,
      name        VARCHAR(100) NOT NULL,
      bio         TEXT,
      onTour      BOOLEAN,
      picture     VARCHAR(1023),
      idLastfm    VARCHAR(100),
      urlLastfm   VARCHAR(1023),
      idDiscogs   VARCHAR(100),
      urlDiscogs  VARCHAR(1023),
      genre       VARCHAR(31),

      artistHash  BIGINT
    );
    

CREATE TABLE IF NOT EXISTS albums(
      id          SERIAL   PRIMARY KEY,
      name        VARCHAR(100) NOT NULL,
      release     VARCHAR(100),
	  urlLastfm   VARCHAR(1023),
	  urlDiscogs  VARCHAR(1023),
      picture     VARCHAR(1023),
      trackCount  INT,

      artistHash  BIGINT,
      albumHash   BIGINT
    );

CREATE TABLE IF NOT EXISTS tracks(
      id          SERIAL   PRIMARY KEY,
	  name        VARCHAR(100) NOT NULL,
      urlLastfm   VARCHAR(1023),
      duration    VARCHAR(100),
      position    VARCHAR(100),

      artistHash  BIGINT,
      albumHash   BIGINT
    );
    
CREATE TABLE IF NOT EXISTS users(
      id          SERIAL   PRIMARY KEY,
      username    VARCHAR(100) UNIQUE NOT NULL,
      firstName   VARCHAR(100) NOT NULL,
      сhatID      INTEGER NOT NULL
    );

CREATE TABLE IF NOT EXISTS user_artist ( -- favourite artists storage --
    id        SERIAL PRIMARY KEY,
    user_id   INTEGER NOT NULL REFERENCES users,
    artist_id INTEGER NOT NULL REFERENCES artists,
    UNIQUE (user_id, artist_id)
);

CREATE TABLE IF NOT EXISTS events (
    id        SERIAL PRIMARY KEY,
    title     VARCHAR(100) NOT NULL,
    titleLink VARCHAR(200),
    date      VARCHAR(100),
    time      VARCHAR(100),
    place     VARCHAR(100),
    placeLink VARCHAR(100),
    cost      VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS artist_event ( 
    id        SERIAL PRIMARY KEY,
    artist_id INTEGER NOT NULL REFERENCES artists,
    event_id  INTEGER NOT NULL REFERENCES events,
    UNIQUE (artist_id, event_id)
);

SELECT * FROM artists;
SELECT * FROM albums;
SELECT * FROM tracks;
SELECT * FROM users;
SELECT * FROM events;
SELECT * FROM user_artist;
SELECT * FROM artist_event;