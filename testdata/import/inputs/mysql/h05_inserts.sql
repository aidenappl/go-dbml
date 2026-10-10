CREATE TABLE p (id INT PRIMARY KEY, name VARCHAR(20), score FLOAT, ok BOOL, born DATE);
INSERT INTO p (id, name, score, ok, born) VALUES (1, 'a', 1.5, TRUE, '2020-01-02'), (2, 'b''s', -3, FALSE, NULL);
INSERT INTO p VALUES (3, "dq", 0, 1, '2021-12-31');
INSERT INTO p SET id = 4, name = 'x';
REPLACE INTO p VALUES (5, 'r', 1, 0, NULL);
INSERT IGNORE INTO p (id) VALUES (6) ON DUPLICATE KEY UPDATE name = 'dup';
