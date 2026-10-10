CREATE TABLE `order_items` (
  `order_id` int,
  `product_id` int,
  `quantity` int,
  PRIMARY KEY (`order_id`, `product_id`)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `order_items` (`order_id`, `product_id`, `quantity`)
VALUES
  (1, 100, 2),
  (1, NULL, 3),
  (NULL, 101, 4);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;