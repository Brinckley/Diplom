CREATE TABLE IF NOT EXISTS artists(
      id          SERIAL   PRIMARY KEY,
      aname       VARCHAR(100) NOT NULL,
      bio         TEXT,
      albumNum    INT,
      onTour      BOOLEAN,
      picture     VARCHAR(255),
      urlLastfm   VARCHAR(255),
      urlDiscogs  VARCHAR(255)
--    lastfmJson  jsonb
--    discogs     jsonb
    );
    

CREATE TABLE IF NOT EXISTS albums(
      id          SERIAL   PRIMARY KEY,
      aname       VARCHAR(100) NOT NULL,
      release     DATE,
      picture     VARCHAR(255),
      trackCount  INT,
	urlLastfm   VARCHAR(255),
	urlDiscogs  VARCHAR(255)
--    lastfmJson  jsonb
--    discogs     jsonb
    );

CREATE TABLE IF NOT EXISTS tracks(
      id          SERIAL   PRIMARY KEY,
	tname       VARCHAR(100) NOT NULL,
      release     DATE,
      lyrics      TEXT,
      urlLastfm  VARCHAR(255)
--    lastfmJson  jsonb
--    discogs     jsonb
    );
    
ALTER TABLE artists 
      ADD COLUMN IF NOT EXISTS album_id INT, 
      ADD CONSTRAINT fk_album 
      FOREIGN KEY (album_id) 
      REFERENCES albums(ID)
      ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE artists 
      ADD COLUMN IF NOT EXISTS track_id INT, 
      ADD CONSTRAINT fk_track 
      FOREIGN KEY (track_id) 
      REFERENCES tracks(ID)
      ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE albums
      ADD COLUMN IF NOT EXISTS artist_id INT, 
      ADD CONSTRAINT fk_artist 
      FOREIGN KEY (artist_id) 
      REFERENCES artists(ID)
      ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE albums
      ADD COLUMN IF NOT EXISTS track_id INT, 
      ADD CONSTRAINT fk_track 
      FOREIGN KEY (track_id) 
      REFERENCES tracks(ID)
      ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE tracks
      ADD COLUMN IF NOT EXISTS artist_id INT, 
      ADD CONSTRAINT fk_artist 
      FOREIGN KEY (artist_id) 
      REFERENCES artists(ID)
      ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE tracks
      ADD COLUMN IF NOT EXISTS album_id INT, 
      ADD CONSTRAINT fk_album 
      FOREIGN KEY (album_id) 
      REFERENCES albums(ID)
      ON UPDATE CASCADE ON DELETE CASCADE;

SELECT * FROM artists;
SELECT * FROM albums;
SELECT * FROM tracks;