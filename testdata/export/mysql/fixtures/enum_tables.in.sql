CREATE TABLE `jobs` (
  `id` integer PRIMARY KEY,
  `status` ENUM ('created', 'running', 'done', 'failed', 'wait for validation') COMMENT 'This is a column note'
);

CREATE TABLE `orders` (
  `id` int UNIQUE PRIMARY KEY,
  `status` ENUM ('created', 'pending', 'processing', 'completed'),
  `created_at` varchar(255)
);
