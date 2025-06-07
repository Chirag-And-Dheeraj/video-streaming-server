ALTER TABLE videos
DROP CONSTRAINT IF EXISTS videos_status_check;

ALTER TABLE videos
RENAME COLUMN status TO upload_status;

ALTER TABLE videos
ADD CONSTRAINT videos_upload_status_check
CHECK (upload_status IN (-1, 0, 1));
