CREATE TABLE videos (
	video_id INTEGER PRIMARY KEY AUTOINCREMENET,
	filename TEXT,
	title TEXT,
	description TEXT,
	upload_initiate_time INTEGER,
	upload_status INTEGER,
	upload_end_time INTEGER,
	manifest_url TEXT
)