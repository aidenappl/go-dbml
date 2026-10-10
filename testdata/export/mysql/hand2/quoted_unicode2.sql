CREATE TABLE `Bảng Người Dùng` (
  `mã số` int PRIMARY KEY,
  `名前` varchar(255),
  `emoji 😀` text COMMENT 'emoji 🎉 inside'
);

CREATE TABLE `order items` (
  `user id` int
);

ALTER TABLE `order items` ADD FOREIGN KEY (`user id`) REFERENCES `Bảng Người Dùng` (`mã số`);
