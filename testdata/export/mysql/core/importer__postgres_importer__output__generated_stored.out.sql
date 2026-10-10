CREATE TABLE `users` (
  `identity_data` jsonb,
  `email` text
);

CREATE TABLE `json_operators` (
  `id` serial PRIMARY KEY AUTO_INCREMENT,
  `data` jsonb,
  `field_text` text,
  `field_json` jsonb,
  `nested_text` text,
  `nested_json` jsonb
);

CREATE TABLE `other_operators` (
  `id` serial PRIMARY KEY AUTO_INCREMENT,
  `a` int,
  `b` int,
  `left_shift` int,
  `right_shift` int,
  `gte` boolean,
  `lte` boolean,
  `neq` boolean,
  `as_text` text
);
