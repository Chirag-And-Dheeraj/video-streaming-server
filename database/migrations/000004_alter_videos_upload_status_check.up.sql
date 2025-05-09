-- Allow -1 in upload_status values
-- Step 1: Drop existing constraint (if it exists)
ALTER TABLE videos
DROP CONSTRAINT IF EXISTS videos_upload_status_check;

-- Step 2: Add new constraint to allow -1, 0, 1
ALTER TABLE videos
ADD CONSTRAINT videos_upload_status_check
CHECK (upload_status IN (-1, 0, 1));
