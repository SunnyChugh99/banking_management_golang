-- Creating "accounts" table
CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Creating "entries" table
CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "amount" bigint NOT NULL
);

-- Creating "transfers" table
CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "amount" bigint NOT NULL
);

-- Creating indexes
CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

-- Adding comments to columns
COMMENT ON COLUMN "entries"."account_id" IS 'Reference to the "id" column in "accounts" table';
COMMENT ON COLUMN "transfers"."amount" IS 'Amount involved in the transfer (must be positive)';

-- Adding foreign key constraints
ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
