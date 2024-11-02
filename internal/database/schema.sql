CREATE TABLE LocalUser (
    id INT PRIMARY KEY,
    name VARCHAR(255),
    pubX CHAR(32),
    pubY CHAR(32),
    privD CHAR(32),
    address VARCHAR(255)
);

CREATE TABLE Posts(
    id INT PRIMARY KEY,
    title VARCHAR(255),
    body VARCHAR(1024),
    imgmd5 VARCHAR(16),
    imgname VAR(255)
);