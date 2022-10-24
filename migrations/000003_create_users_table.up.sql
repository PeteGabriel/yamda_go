CREATE TABLE IF NOT EXISTS users (
  id bigint(20) PRIMARY KEY,
  created_at timestamp(0) NOT NULL DEFAULT NOW(),
  name varchar(200) NOT NULL,
  email varchar(50) UNIQUE NOT NULL,
  password_hash BINARY(20) NOT NULL,
  activated bool NOT NULL,
  version int NOT NULL DEFAULT 1
);
