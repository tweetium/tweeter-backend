CREATE OR REPLACE FUNCTION update_created_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.created = now();
  NEW.modified = now();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_createtime BEFORE INSERT ON users FOR EACH ROW EXECUTE PROCEDURE update_created_column();
