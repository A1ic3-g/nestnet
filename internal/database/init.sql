CREATE TABLE LocalUser (
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    pubX CHAR(32) NOT NULL,
    pubY CHAR(32) NOT NULL,
    privD CHAR(32) NOT NULL,
    address VARCHAR(255) NOT NULL
);

CREATE TABLE Posts(
    id INT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    body VARCHAR(1024) NOT NULL,
    imgmd5 VARCHAR(16),
    imgname VARCHAR(255)
);