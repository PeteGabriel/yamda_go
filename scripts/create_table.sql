CREATE TABLE IF NOT EXISTS Movie (
    `Id` bigint(20) NOT NULL AUTO_INCREMENT,
    `created_at` timestamp(0) NOT NULL DEFAULT NOW(),
    `title` varchar(200) NOT NULL,
    `year` int NOT NULL,
    `runtime` int NOT NULL,
    `genres` TINYTEXT NOT NULL,
    `version` int NOT NULL DEFAULT 1,
    PRIMARY KEY (`Id`),
    UNIQUE KEY `ID_UNIQUE` (`Id`)
    );