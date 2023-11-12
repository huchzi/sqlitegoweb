DROP TABLE IF EXISTS album;
CREATE TABLE album (
  id         INTEGER PRIMARY KEY,
  title      VARCHAR(128) NOT NULL,
  artist     VARCHAR(255) NOT NULL,
  price      DECIMAL(5,2) NOT NULL
);

INSERT INTO album
  (id, title, artist, price)
VALUES
  (1, 'Blue Train', 'John Coltrane', 56.99),
  (2, 'Giant Steps', 'John Coltrane', 63.99),
  (3, 'Jeru', 'Gerry Mulligan', 17.99),
  (4, 'Sarah Vaughan', 'Sarah Vaughan', 34.98);