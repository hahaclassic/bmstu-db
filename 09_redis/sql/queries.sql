CREATE TABLE IF NOT EXISTS temp_tracks (
    id UUID,
    name VARCHAR(100),
    album VARCHAR(100),
    explicit BOOLEAN,
    duration INT,
    genre VARCHAR(50),
    stream_count BIGINT
);

drop table temp_tracks;

INSERT INTO temp_tracks (id, name, album, explicit, duration, genre, stream_count)
SELECT 
    t.id, 
    t.name, 
    (select a.title from albums a where a.id = t.album_id) as album,
    t.explicit, 
    t.duration, 
    t.genre, 
    t.stream_count
FROM 
    tracks t
JOIN 
    albums a 
ON 
    t.album_id = a.id;