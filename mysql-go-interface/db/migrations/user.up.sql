-- +migrate Up
CREATE TABLE IF NOT EXISTS user_settings(
   id binary(16) PRIMARY KEY,
   two_steps_verif BOOL DEFAULT false, 
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS users(
   id binary(16) PRIMARY KEY,
   name VARCHAR (50) UNIQUE NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   password BINARY (60) NOT NULL,
   user_settings_id binary(16),
   FOREIGN KEY (user_settings_id) REFERENCES user_settings(id), 
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);
