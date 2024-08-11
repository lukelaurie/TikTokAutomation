CREATE TABLE preferences (
    id SERIAL PRIMARY KEY,
    username TEXT REFERENCES users(username),
    video_type TEXT NOT NULL,
    background_video_type TEXT NOT NULL,
    font_name TEXT NOT NULL,
    font_color TEXT NOT NULL,
    preference_order INT
);