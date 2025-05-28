ALTER TABLE videos
DROP CONSTRAINT IF EXISTS videos_upload_status_check;

ALTER TABLE videos
RENAME COLUMN upload_status TO status;

ALTER TABLE videos
ADD CONSTRAINT videos_status_check
CHECK (status IN (-1, 0, 1, 2));
