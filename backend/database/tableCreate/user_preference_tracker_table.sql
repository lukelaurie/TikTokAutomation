CREATE TABLE user_preference_tracker (
    id SERIAL PRIMARY KEY,
    username TEXT REFERENCES users(username),
    current_preference_order INT DEFAULT 1,
    current_preference_count INT DEFAULT 0
);