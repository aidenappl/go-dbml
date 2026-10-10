CREATE TABLE `A` (
  `id` integer,
  `code` number
);

CREATE TABLE `B` (
  `id` integer,
  `code` number
);

ALTER TABLE `A` ADD FOREIGN KEY (`id`) REFERENCES `B` (`id`) ON DELETE NO ACTION ON UPDATE CASCADE;

ALTER TABLE `A` ADD FOREIGN KEY (`code`) REFERENCES `B` (`code`) ON DELETE SET DEFAULT ON UPDATE SET NULL;
