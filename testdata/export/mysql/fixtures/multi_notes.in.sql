CREATE TABLE `orders` (
  `id` int PRIMARY KEY COMMENT 'primary key field',
  `status` ENUM ('created', 'pending', 'processing', 'completed'),
  `created_at` varchar(255)
);

CREATE TABLE `bookings` (
  `id` integer,
  `country` varchar(255),
  `booking_date` date,
  `created_at` timestamp,
  PRIMARY KEY (`id`, `country`)
);

CREATE INDEX `created_at_index` ON `bookings` (`created_at`);

CREATE INDEX `bookings_index_1` ON `bookings` (`booking_date`);

CREATE UNIQUE INDEX `bookings_index_2` ON `bookings` (`country`, `booking_date`);

CREATE INDEX `bookings_index_3` ON `bookings` (`booking_date`) USING HASH;

CREATE INDEX `bookings_index_4` ON `bookings` ((id*2));

CREATE INDEX `bookings_index_5` ON `bookings` ((id*3), (getdate()));

CREATE INDEX `bookings_index_6` ON `bookings` ((id*3), `id`);

ALTER TABLE `orders` COMMENT = 'Note on table orders';
