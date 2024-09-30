COPY artists (id, name, genre, country, debut_year)
FROM '/data/artists.csv' DELIMITER ',' CSV HEADER;

COPY albums (id, title, release_date, label, genre)
FROM '/data/albums.csv' DELIMITER ',' CSV HEADER;

COPY tracks (id, name, order_in_album, album_id, explicit, duration, genre, stream_count)
FROM '/data/tracks.csv' DELIMITER ',' CSV HEADER;

COPY playlists (id, title, description, private, last_updated, rating)
FROM '/data/playlists.csv' DELIMITER ',' CSV HEADER;

COPY users (id, name, registration_date, birth_date, premium, premium_expiration)
FROM '/data/users.csv' DELIMITER ',' CSV HEADER;

COPY playlist_tracks (track_id, playlist_id, date_added, track_order)
FROM '/data/playlist_tracks.csv' DELIMITER ',' CSV HEADER;

COPY user_playlists (playlist_id, user_id, is_favorite, access_level)
FROM '/data/user_playlists.csv' DELIMITER ',' CSV HEADER;

COPY tracks_by_artists (track_id, artist_id)
FROM '/data/tracks_by_artists.csv' DELIMITER ',' CSV HEADER;

COPY albums_by_artists (album_id, artist_id)
FROM '/data/albums_by_artists.csv' DELIMITER ',' CSV HEADER;
