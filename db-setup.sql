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
    price_per_liter_euro    float8 constraint price_euro_not_negative_or_zero CHECK (price_per_liter_euro IS NULL OR price_per_liter_euro > 0),
    total_liter             float8 NOT NULL constraint total_liter_not_negative_or_zero CHECK (total_liter > 0),
    price_per_liter         float8 constraint price_not_negative_or_zero CHECK (price_per_liter IS NULL OR price_per_liter > 0),
    currency                varchar(12),
    distance                float8 NOT NULL constraint distance_not_negative_or_zero CHECK (distance > 0),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON refuel
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

INSERT INTO public.refuel ("name",date_time,price_per_liter_euro,total_liter,price_per_liter,currency,distance) VALUES
	 ('LPG Kreuzlingen','2021-09-04 13:10:25',1.439,42.0,1.488,'chf', 560.6);
INSERT INTO public.refuel ("name",date_time,price_per_liter_euro,total_liter,price_per_liter,currency,distance) VALUES
	 ('LPG Kreuzlingen Test','2021-09-15 16:12:11',1.368,34.0,1.420,'chf', 420.8);