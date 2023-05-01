CREATE TABLE users (
  id TEXT PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255) UNIQUE,
  password VARCHAR
);

CREATE TABLE categories (
  id TEXT PRIMARY KEY,
  name VARCHAR(255),
  description VARCHAR(255)
);
