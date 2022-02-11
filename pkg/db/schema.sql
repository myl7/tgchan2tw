CREATE TABLE IF NOT EXISTS tg_in (
  id TEXT PRIMARY KEY
);
CREATE TABLE IF NOT EXISTS tw_out (
  id INTEGER(8) PRIMARY KEY
);
CREATE TABLE IF NOT EXISTS tg_in_to_tw_out (
  tg_in_id TEXT REFERENCES tg_in (id),
  tw_out_id INTEGER(8) REFERENCES tw_out (id)
);
