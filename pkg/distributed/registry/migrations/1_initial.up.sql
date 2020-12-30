CREATE TABLE segments (
  id varchar(255) PRIMARY KEY UNIQUE, 
  order_id varchar(255), 
  input_storage_claim_identity TEXT,
  output_storage_claim_identity TEXT,
  kind varchar(255), 
  locked_until varchar(255), 
  locked_by varchar(255), 
  created_at varchar(255), 
  updated_at varchar(255), 
  payload TEXT
);

CREATE TABLE orders (
  id varchar(255) PRIMARY KEY UNIQUE, 
  kind varchar(255), 
  payload TEXT
);
