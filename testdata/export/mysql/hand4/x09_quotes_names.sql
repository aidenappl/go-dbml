CREATE TABLE `Order Items` (
  `Item ID` int PRIMARY KEY,
  `Type` ENUM ('a b', 'c'),
  `user` USER,
  `qty` Int,
  `select` varchar(255)
);
