CREATE TABLE `categories` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `products` (
  `id` int PRIMARY KEY,
  `category_id` int,
  `name` varchar(255)
);

ALTER TABLE `products` ADD FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `categories` (`id`, `name`)
VALUES
  (1, 'Electronics');
INSERT INTO `products` (`id`, `category_id`, `name`)
VALUES
  (1, 1, 'Laptop'),
  (2, NULL, 'Uncategorized Item');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;