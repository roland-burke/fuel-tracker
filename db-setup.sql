CREATE DATABASE fuel_tracker OWNER postgres;

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS refuel (
    id                      SERIAL NOT NULL PRIMARY key,
    name                    varchar(50) NOT NULL,
    date_time               timestamp NOT NULL,
    price_per_liter_euro    float8,
    total_liter             float8 NOT NULL,
    price_per_liter         float8,
    currency                varchar(15),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON refuel
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

INSERT INTO public.refuel ("name",date_time,price_per_liter_euro,total_liter,price_per_liter,currency) VALUES
	 ('LPG Kreuzlingen','2021-09-04 13:10:25',1.439,42.0,1.488,'chf');