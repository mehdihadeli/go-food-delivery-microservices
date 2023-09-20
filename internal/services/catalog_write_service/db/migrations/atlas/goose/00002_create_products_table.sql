CREATE TABLE IF NOT EXISTS products
(
    product_id  uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        text,
    description text,
    price       numeric,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone
);
