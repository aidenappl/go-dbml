CREATE TABLE `test` (
  `text1` varchar(255) DEFAULT 'Test2',
  `text2` varchar(255) DEFAULT 'Test2',
  `text3` varchar(255) DEFAULT 'Test2',
  `text4` varchar(255) DEFAULT 'Test2',
  `number1` int DEFAULT 1,
  `number2` int DEFAULT 1,
  `number3` int DEFAULT 1,
  `number4` int DEFAULT 1,
  `date1` datetime DEFAULT (GETDATE()),
  `date2` datetime DEFAULT (GETDATE()),
  `date3` datetime DEFAULT (GETDATE()),
  `date4` datetime DEFAULT (GETDATE())
);
