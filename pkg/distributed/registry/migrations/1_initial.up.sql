CREATE TABLE segments (
  id varchar(255) PRIMARY KEY UNIQUE, 
  order_id varchar(255), 
  storage_claim_identity TEXT,
  kind varchar(255), 
  payload TEXT
);

CREATE TABLE orders (
  id varchar(255) PRIMARY KEY UNIQUE, 
  kind varchar(255), 
  payload TEXT
);
