DROP DATABASE IF EXISTS balanceApp;

CREATE DATABASE balanceApp;

USE balanceApp;

CREATE TABLE `accounts` (
    `userID` int NOT NULL PRIMARY KEY,
    `amount` float
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `orders` (
    `id` int NOT NULL PRIMARY KEY NOT NULL,
    `orderID` int NOT NULL,
    `userID` int NOT NULL,
    `serviceType` int NOT NULL,
    `orderCost` float,
    `creatingTime` datetime,
    `comments` text,
    `orderState` int,
    FOREIGN KEY (`userID`)  REFERENCES `accounts` (`userID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `transactions` (
    `id` int NOT NULL PRIMARY KEY NOT NULL,
    `transactionID` int NOT NULL,
    `userID` int NOT NULL,
    `transactionType` int,
    `sum` float,
    `time` datetime,
    `actionComments` text,
    `addComments` text,
    FOREIGN KEY (`userID`)  REFERENCES `accounts` (`userID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
