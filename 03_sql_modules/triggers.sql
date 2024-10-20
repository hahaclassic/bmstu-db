CREATE OR REPLACE FUNCTION update_last_updated() 
RETURNS TRIGGER AS $$
BEGIN
    UPDATE playlists
    SET last_updated = CURRENT_TIMESTAMP
    WHERE id = NEW.playlist_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql; 

CREATE TRIGGER update_playlist_timestamp
AFTER INSERT OR UPDATE ON playlist_tracks
FOR EACH ROW
EXECUTE FUNCTION update_last_updated();