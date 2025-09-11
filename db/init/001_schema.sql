CREATE
DATABASE IF NOT EXISTS app CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE
app;
CREATE TABLE IF NOT EXISTS users
(
    id
    VARCHAR
(
    64
) PRIMARY KEY,
    name VARCHAR
(
    255
) NOT NULL,
    created_at DATETIME
(
    6
) NOT NULL,
    INDEX idx_created_at
(
    created_at
    DESC
)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE =utf8mb4_0900_ai_ci;