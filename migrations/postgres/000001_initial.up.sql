CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_type AS ENUM (
    '1fe92aa8-2a61-4bf1-b907-182b497584ad', -- system user
    '9fb3ada6-a73b-4b81-9295-5c1605e54552'  -- admin user
);

CREATE TYPE app_type AS ENUM (
    '1fe92aa8-2a61-4bf1-b907-182b497584ad', -- client
    '9fb3ada6-a73b-4b81-9295-5c1605e54552'  -- admin
);

CREATE TABLE IF NOT EXISTS "user" (
    "id" UUID PRIMARY KEY,
    "user_type_id" user_type NOT NULL,
    "first_name" VARCHAR(250) NOT NULL,
    "last_name" VARCHAR(250) NOT NULL,
    "phone_number" VARCHAR(30) NOT NULL,
    "image" TEXT,
    "deleted_at" BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX "user_deleted_at_idx" ON "user"("deleted_at");

INSERT INTO "user" (
    "id",
    "first_name",
    "last_name",
    "phone_number",
    "user_type_id"
) VALUES (
    '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd',
    'admin',
    'admin',
    '99894172774',
    '9fb3ada6-a73b-4b81-9295-5c1605e54552'
);


CREATE TABLE "company_type" (
  "id" UUID PRIMARY KEY,
  "name" JSONB NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" UUID,
  "deleted_at" BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX "company_type_deleted_at_idx" ON "company_type"("deleted_at");

INSERT INTO  "company_type" (
  "id",
  "name",
  "created_by"
) VALUES 
  ('ab7c6a4f-4fc9-486c-86cf-13bf1ced9284', '{"en":"Warehouse", "uz":"Ombor", "ru":"Склад"}','9a2aa8fe-806e-44d7-8c9d-575fa67ebefd' ),
  ('2321ec59-e85e-4fef-bca3-9e6c87a2cee3', '{"en":"Grocery", "uz":"Oziq-ovqat", "ru":"Продуктовый"}','9a2aa8fe-806e-44d7-8c9d-575fa67ebefd');


CREATE TABLE "company_size" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR(200) NOT NULL,
  "name_tr" JSONB NOT NULL,
  "from" INT NOT NULL,
  "to" INT NOT NULL,
  "description" JSONB,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" UUID,
  "deleted_at" BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX "company_size_deleted_at_idx" ON "company_size"("deleted_at");

INSERT INTO  "company_size" (
  "id",
  "name",
  "name_tr",
  "from",
  "to",
  "created_by"
) VALUES 
  ('ab7c6a4f-4fc9-486c-86cf-13bf1ced9284', 'Small', '{"en":"Small", "uz":"Kichik", "ru":"Maленький"}', 0, 20, '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd'),
  ('2321ec59-e85e-4fef-bca3-9e6c87a2cee3', 'Medium', '{"en":"Medium", "uz":"O''rtacha", "ru":"Средний"}', 20, 100,'9a2aa8fe-806e-44d7-8c9d-575fa67ebefd');




CREATE TABLE "company" (
  "id" UUID PRIMARY KEY,
  "name" varchar(200) NOT NULL,
  "owner" UUID,
  "email" VARCHAR,
  "legal_name" VARCHAR,
  "legal_adress" VARCHAR,
  "country" VARCHAR,
  "zip_code" VARCHAR,
  "tax_payer_id" VARCHAR,
  "type_id" UUID NOT NULL REFERENCES "company_type"("id"),
  "size_id" UUID NOT NULL REFERENCES "company_size"("id"),
  "created_by" UUID,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID,
  UNIQUE ("name", "deleted_at")
);

CREATE INDEX "company_deleted_at_idx" ON "company"("deleted_at");

CREATE TABLE "company_user" (
  "user_id" UUID NOT NULL,
  "company_id" UUID NOT NULL,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY ("user_id", "company_id", "deleted_at")
);
CREATE INDEX "company_user_deleted_at" ON "company_user"("deleted_at");

CREATE TABLE "shop" (
  "id" UUID PRIMARY KEY ,
  "title" VARCHAR NOT NULL,
  "phone_number" VARCHAR NOT NULL,
  "company_id" UUID,
  "size" INT NOT NULL DEFAULT 0,
  "address" VARCHAR NOT NULL DEFAULT '',
  "description" TEXT NOT NULL DEFAULT '',
  "number_of_cashboxes" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "created_by" UUID,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID,
  UNIQUE ("title", "company_id", "deleted_at")
);
CREATE INDEX "shop_deleted_at_idx" ON "shop"("deleted_at");


CREATE TABLE "cheque" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "name" VARCHAR(300) NOT NULL,
  "message" TEXT  NOT NULL,
  "created_by" UUID,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID,
  UNIQUE ("company_id", "name", "deleted_at")
);

CREATE INDEX "cheque_deleted_at_idx" ON "cheque"("deleted_at"); 

CREATE TABLE "cashbox" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "shop_id" UUID,
  "title" VARCHAR(200) NOT NULL,
  "cheque_id" UUID NOT NULL REFERENCES "cheque"("id") ON DELETE CASCADE,
  "created_by" UUID,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID,
  UNIQUE ("company_id", "title", "deleted_at")
);
CREATE INDEX "cashbox_deleted_at" ON "cashbox"("deleted_at");


CREATE TABLE "receipt_block" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR(200) NOT NULL,
  "name_tr" JSONB NOT NULL,
  "created_by" UUID,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID
);

CREATE INDEX "receipt_block_deleted_at_idx" ON "receipt_block"("deleted_at");

-- insert receipt block
INSERT INTO "receipt_block" ("id", "name", "name_tr", "created_by") VALUES
    ('9bdde8d7-1b48-4788-9f6f-aa94bd6d9006', 'information block', '{"en":"Information Block", "uz":"Axborot bloki", "ru":"Информационный блок"}', '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd'),
    ('a7655040-2942-4b61-9771-847f4f48a33d', 'Bottom block', '{"en":"Bottom Block", "uz":"Pastki blok", "ru":"Нижний блок"}', '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd');


CREATE TABLE "receipt_field"  (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR(200) NOT NULL,
  "name_tr" JSONB NOT NULL,
  "block_id" UUID NOT NULL REFERENCES "receipt_block"("id") ON DELETE CASCADE,
  "created_by" UUID REFERENCES "user"("id") ON DELETE SET NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID REFERENCES "user"("id") ON DELETE SET NULL,
  UNIQUE ("name", "deleted_at")
);

CREATE INDEX "receipt_field_deleted_at_idx" ON "receipt_field"("deleted_at");

INSERT INTO "receipt_field"("id", "name", "name_tr", "block_id", "created_by" ) VALUES
    ('9527f6aa-6432-4932-ab43-5bfed278e751', 'shop name', '{"en":"Shop Name", "ru":"Название магазина", "uz":"Do''kon nomi"}', '9bdde8d7-1b48-4788-9f6f-aa94bd6d9006', '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd'),
    ('bc289940-a175-4460-a743-0aa3da9f8c7c', 'datetime', '{"en":"Datetime", "ru":"Дата и время", "uz":"Sana vaqti"}', '9bdde8d7-1b48-4788-9f6f-aa94bd6d9006', '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd'),
    ('1888a21b-110a-41a6-8c97-5d7b8f34ea3b', 'seller', '{"en":"Seller", "ru":"Продавец", "uz":"Sotuvchi"}', '9bdde8d7-1b48-4788-9f6f-aa94bd6d9006', '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd'),
    ('c2e9abae-1c25-4ac1-85c3-4c8c33d0d975', 'cashier', '{"en":"Cashier", "ru":"Кассир", "uz":"Kassir"}', '9bdde8d7-1b48-4788-9f6f-aa94bd6d9006', '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd'),
    ('15f8c6fe-8ef2-4eb2-bcac-1cf52402a4cd', 'customer', '{"en":"Customer", "ru":"Клиент", "uz":"Mijoz"}', '9bdde8d7-1b48-4788-9f6f-aa94bd6d9006', '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd'),
    ('ab9c342d-8d6f-4e82-9622-11437c6e135c', 'contacts', '{"en":"Contacts", "ru":"Контакты", "uz":"Kontaktlar"}', '9bdde8d7-1b48-4788-9f6f-aa94bd6d9006', '9a2aa8fe-806e-44d7-8c9d-575fa67ebefd');


CREATE  TABLE "cheque_logo" (
  "image" TEXT NOT NULL,
  "cheque_id" UUID NOT NULL REFERENCES "cheque"("id") ON DELETE CASCADE,
  "left" INT NOT NULL DEFAULT 0,
  "right" INT NOT NULL DEFAULT 0,
  "top" INT NOT NULL DEFAULT 0,
  "bottom" INT NOT NULL DEFAULT 0,
  PRIMARY KEY("cheque_id")
);

CREATE TABLE "cheque_field" (
  "field_id" UUID NOT NULL REFERENCES "receipt_field"("id") ON DELETE CASCADE,
  "cheque_id" UUID NOT NULL REFERENCES "cheque"("id") ON DELETE CASCADE,
  "position" INT NOT NULL,
  "is_added" BOOLEAN DEFAULT FALSE,
  "created_by" UUID REFERENCES "user"("id") ON DELETE SET NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID REFERENCES "user"("id") ON DELETE SET NULL,
  PRIMARY KEY ("cheque_id", "field_id", "position", "deleted_at")
);
CREATE INDEX "cheque_field_deleted_at" ON "cheque_field"("deleted_at");

CREATE TABLE "payment_type" (
  "id" UUID PRIMARY KEY,
  "name" varchar NOT NULL,
  "logo" TEXT,
  "company_id" UUID NOT NULL REFERENCES "company"("id") ON DELETE CASCADE,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" UUID REFERENCES "user"("id") ON DELETE SET NULL,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID REFERENCES "user"("id") ON DELETE SET NULL,
  UNIQUE ("name", "company_id", "deleted_at")
);

CREATE INDEX "payment_type_deleted_at_idx" ON "payment_type"("deleted_at");


CREATE TABLE "cashbox_payment" (
  "id" UUID PRIMARY KEY,
  "cashbox_id" UUID NOT NULL REFERENCES "cashbox"("id") ON DELETE CASCADE,
  "payment_type_id" UUID NOT NULL REFERENCES "payment_type"("id") ON DELETE CASCADE,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" UUID REFERENCES "user"("id") ON DELETE SET NULL,
  "deleted_at" BIGINT NOT NULL DEFAULT 0,
  "deleted_by" UUID REFERENCES "user"("id") ON DELETE SET NULL,
  UNIQUE ("cashbox_id", "payment_type_id", "deleted_at")
);

CREATE INDEX "cashbox_payment_deleted_at_idx" ON "cashbox_payment"("deleted_at");

-- trigger
CREATE OR REPLACE FUNCTION cashbox_create()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
DECLARE
  number_of_cashbox INT := 0;
BEGIN
  SELECT count(1) FROM "cashbox" WHERE "shop_id"=NEW.shop_id AND "deleted_at" = 0 INTO number_of_cashbox;

  UPDATE shop SET "number_of_cashboxes" = number_of_cashbox WHERE "id"=NEW."shop_id" AND deleted_at = 0;

  RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER cashbox_create_or_update
  AFTER INSERT OR UPDATE ON "cashbox"
  FOR EACH ROW
  EXECUTE PROCEDURE cashbox_create();

-- trigger

CREATE OR REPLACE FUNCTION create_defaults()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
DECLARE
  shop_id UUID := uuid_generate_v4();
  user_phone_number VARCHAR := '';
  cheque_id UUID := uuid_generate_v4();
BEGIN
    SELECT "phone_number" FROM "user" WHERE "id" = NEW.created_by INTO user_phone_number;

    --  shop 
    INSERT INTO "shop" ("id", "company_id", "title","address", "phone_number", "created_by") VALUES (shop_id,  NEW."id", CONCAT('Store ', NEW.name), '',  user_phone_number, NEW.created_by );

    -- standart cheque
    INSERT INTO "cheque" ("id", "name", "company_id", "message", "created_by") VALUES (cheque_id, 'Standart', NEW."id", 'Thank you for your purchase!', NEW."created_by");
    -- cheque fields
    INSERT INTO "cheque_field" ("field_id", "cheque_id", "position", "created_by") VALUES 
      ('9527f6aa-6432-4932-ab43-5bfed278e751', cheque_id, 1, NEW.created_by),
      ('bc289940-a175-4460-a743-0aa3da9f8c7c', cheque_id, 2, NEW.created_by),
      ('1888a21b-110a-41a6-8c97-5d7b8f34ea3b', cheque_id, 3, NEW.created_by),
      ('c2e9abae-1c25-4ac1-85c3-4c8c33d0d975', cheque_id, 4, NEW.created_by),
      ('15f8c6fe-8ef2-4eb2-bcac-1cf52402a4cd', cheque_id, 5, NEW.created_by),
      ('ab9c342d-8d6f-4e82-9622-11437c6e135c', cheque_id, 6, NEW.created_by);

    -- payment type
    INSERT INTO "payment_type" ("id", "company_id", "name", "created_by") VALUES
        (uuid_generate_v4(), NEW."id", 'CASH', NEW.created_by),
        (uuid_generate_v4(), NEW."id", 'UZCARD', NEW.created_by),
        (uuid_generate_v4(), NEW."id", 'HUMO', NEW.created_by),
        (uuid_generate_v4(), NEW."id", 'VISA', NEW.created_by),
        (uuid_generate_v4(), NEW."id", 'MASTERCARD', NEW.created_by);
    -- cashbox 
    INSERT INTO "cashbox" ("id", "title", "shop_id", "company_id", "created_by", "cheque_id") VALUES (uuid_generate_v4(), 'Cashbox', "shop_id", NEW."id", NEW.created_by, cheque_id );
	RETURN NEW;
END;
$$;


-- triggers
CREATE OR REPLACE TRIGGER create_defaults
    AFTER INSERT ON "company"
    FOR EACH ROW
    EXECUTE PROCEDURE create_defaults();

