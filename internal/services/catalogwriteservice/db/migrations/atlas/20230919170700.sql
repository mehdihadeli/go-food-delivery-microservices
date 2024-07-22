-- Create "products" table
CREATE TABLE "public"."products" (
  "id" text NOT NULL,
  "name" text NULL,
  "description" text NULL,
  "price" numeric NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
