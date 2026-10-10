CREATE TABLE `orders` (
  `id` int PRIMARY KEY,
  `status` ENUM ('pending', 'processing', 'completed') DEFAULT 'pending'
);
