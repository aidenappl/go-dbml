CREATE TABLE `d` (
  `a` int,
  `b` int,
  `c` bool DEFAULT true,
  `e` bool DEFAULT false,
  `f` varchar(255) DEFAULT '',
  `g` int DEFAULT -1.25,
  `h` text DEFAULT ('quoted'),
  `i` int DEFAULT 1000,
  `j` int UNIQUE PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `k` varchar(255)
);
