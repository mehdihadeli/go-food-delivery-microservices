-- Create "products" table
CREATE TABLE "public"."products" (
  "product_id" text NOT NULL,
  "name" text NULL,
  "description" text NULL,
  "price" numeric NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("product_id")
);
