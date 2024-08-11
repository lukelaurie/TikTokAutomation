CREATE TABLE public."users"
(
	username TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
);