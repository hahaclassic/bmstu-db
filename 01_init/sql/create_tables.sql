CREATE extension IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS artists (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    genre VARCHAR(50),
    country VARCHAR(50),
    debut_year INT
)

CREATE TABLE IF NOT EXISTS albums (
    id UUID PRIMARY KEY,
    title VARCHAR(100),
    release_date DATE,
    label VARCHAR(100),
    genre VARCHAR(50)
)

CREATE TABLE IF NOT EXISTS tracks (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    explicit BOOLEAN,
    duration INT,
    genre VARCHAR(50),
    stream_count BIGINT,
)

CREATE TABLE IF NOT EXISTS playlists (
    id UUID PRIMARY KEY,
    title VARCHAR(100),
    description TEXT,
    private BOOLEAN,
    last_updated TIMESTAMP,
    rating INT
)

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    registration_date TIMESTAMP,
    birth_date DATE,
    premium BOOLEAN,
    premium_expiration TIMESTAMP,
)

CREATE TABLE IF NOT EXISTS playlist_tracks (
    track_id UUID,
    playlist_id UUID,
    date_added TIMESTAMP,
    track_order INT
)

CREATE TABLE IF NOT EXISTS album_tracks (
    track_id UUID,
    album_id UUID,
    track_order INT
);

CREATE TABLE IF NOT EXISTS user_playlists (
    playlist_id UUID,
    user_id UUID,
    is_favorite BOOLEAN,
    access_level INT
)

CREATE TABLE IF NOT EXISTS albums_by_artists (
    album_id UUID,
    artist_id UUID,
);

CREATE TABLE IF NOT EXISTS tracks_by_artists (
    track_id UUID,
    artist_id UUID
);