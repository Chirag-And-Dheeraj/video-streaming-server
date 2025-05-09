-- Revert to only allowing 0 and 1 in upload_status
-- Step 1: Drop updated constraint
ALTER TABLE videos
DROP CONSTRAINT IF EXISTS videos_upload_status_check;

-- Step 2: Re-add original constraint to only allow 0 and 1
ALTER TABLE videos
ADD CONSTRAINT videos_upload_status_check
CHECK (upload_status IN (0, 1));
