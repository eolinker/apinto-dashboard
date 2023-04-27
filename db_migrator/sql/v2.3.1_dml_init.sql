INSERT INTO `user`(`sex`, `username`,`password`)

SELECT * FROM (SELECT 0, "admin","550e1bafe077ff0b0b67f4e32f29d751") A

WHERE NOT EXISTS (SELECT NULL FROM `user`);