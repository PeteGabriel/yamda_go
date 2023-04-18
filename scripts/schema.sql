CREATE TABLE IF NOT EXISTS Movie (
    `Id` bigint(20) NOT NULL AUTO_INCREMENT,
    `created_at` timestamp(0) NOT NULL DEFAULT NOW(),
    `title` varchar(200) NOT NULL,
    `year` int NOT NULL,
    `runtime` int NOT NULL,
    `genres` TINYTEXT NOT NULL,
    `version` int NOT NULL DEFAULT 1,
    PRIMARY KEY (`Id`),
    UNIQUE KEY `ID_UNIQUE` (`Id`),
    FULLTEXT(`genres`)
    );

ALTER TABLE Movie ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);

INSERT INTO Movie (title, year, runtime, genres) VALUES ("Gothika", 2003, 125, "Horror,Thriller");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("Broken Embraces", 2009, 135, "Drama,Romance,Thriller");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("A Christmas Prince: The Royal Baby", 2019, 155, "Romance,Family");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("Deconstructing Harry", 1997, 115, "Comedy,Drama");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("Tie Me Up! Tie Me Down!", 1989, 203, "Crime,Comedy,Drama");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("Grudge Match", 2013, 115, "Comedy");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("Halloween II", 1981, 122, "Horror,Thriller");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("Bill & Ted Face the Music", 2020, 175, "Science Fiction,Adventure,Comedy");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("From Dusk Till Dawn", 1996, 145, "Action,Thriller,Crime");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("The Knight Before Christmas", 1997, 115, "Comedy,Romance");
INSERT INTO Movie (title, year, runtime, genres) VALUES ("The Ghost and the Darkness", 1996, 165, "Adventure");



CREATE TABLE IF NOT EXISTS users (
  id bigint(20) PRIMARY KEY AUTO_INCREMENT,
  created_at timestamp(0) NOT NULL DEFAULT NOW(),
  name varchar(200) NOT NULL,
  email varchar(50) UNIQUE NOT NULL,
  password_hash BINARY(28) NOT NULL,
  activated bool NOT NULL,
  version int NOT NULL DEFAULT 1
  );
