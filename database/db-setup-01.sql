CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
    users_id                SERIAL NOT NULL PRIMARY key,
    username                varchar(30) NOT NULL,
    pass_key                varchar(30) NOT NULL,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refuel (
    id                      SERIAL NOT NULL PRIMARY key,
    users_id                int,
    description             varchar(50) NOT NULL,
    date_time               timestamp NOT NULL,
    price_per_liter_euro    float8 NOT NULL constraint price_euro_not_negative CHECK (price_per_liter_euro >= 0),
    total_liter             float8 NOT NULL constraint total_liter_not_negative_or_zero CHECK (total_liter > 0),
    price_per_liter         float8 NOT NULL constraint price_not_negative CHECK (price_per_liter >= 0),
    currency                varchar(12) NOT NULL,
    mileage                 float8 NOT NULL constraint mileage_not_negative_or_zero CHECK (mileage > 0),
    license_plate           varchar(15) NOT NULL,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    constraint fk_users foreign key(users_id) references users(users_id)
);

CREATE TRIGGER set_timestamp_refuel
BEFORE UPDATE ON refuel
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();