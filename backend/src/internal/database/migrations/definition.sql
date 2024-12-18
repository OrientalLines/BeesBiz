DO $$ BEGIN IF NOT EXISTS (
	SELECT
		1
	FROM
		pg_type typ
		INNER JOIN pg_namespace nsp ON nsp.oid = typ.typnamespace
	WHERE
		nsp.nspname = current_schema()
		AND typ.typname = 'role'
) THEN CREATE TYPE role AS ENUM ('ADMIN', 'WORKER', 'MANAGER');

END IF;

END;

$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS "region" (
	"region_id" SERIAL,
	"name" VARCHAR NOT NULL,
	"climate_zone" VARCHAR,
	PRIMARY KEY("region_id")
);

CREATE TABLE IF NOT EXISTS "apiary" (
	"apiary_id" SERIAL,
	"location" VARCHAR NOT NULL,
	"manager_id" INTEGER,
	"establishment_date" DATE,
	PRIMARY KEY("apiary_id")
);

CREATE TABLE IF NOT EXISTS "hive" (
	"hive_id" SERIAL,
	"apiary_id" INTEGER,
	"hive_type" VARCHAR,
	"installation_date" DATE,
	"current_status" VARCHAR,
	PRIMARY KEY("hive_id")
);

CREATE TABLE IF NOT EXISTS "bee_community" (
	"community_id" SERIAL,
	"hive_id" INTEGER,
	"queen_age" INTEGER,
	"population_estimate" INTEGER,
	"health_status" VARCHAR,
	PRIMARY KEY("community_id")
);

CREATE TABLE IF NOT EXISTS "veterinary_passport" (
	"passport_id" SERIAL,
	"bee_community_id" INTEGER,
	"issue_date" DATE,
	"health_status" VARCHAR,
	"last_inspection_date" DATE,
	PRIMARY KEY("passport_id")
);

CREATE TABLE IF NOT EXISTS "veterinary_record" (
	"record_id" SERIAL,
	"passport_id" INTEGER,
	"record_date" DATE,
	"description" TEXT,
	"treatment" TEXT,
	PRIMARY KEY("record_id")
);

CREATE TABLE IF NOT EXISTS "sensor" (
	"sensor_id" SERIAL,
	"hive_id" INTEGER,
	"sensor_type" VARCHAR,
	"last_reading" BYTEA,
	"last_reading_time" TIMESTAMP,
	PRIMARY KEY("sensor_id")
);

CREATE TABLE IF NOT EXISTS "sensor_reading" (
	"reading_id" SERIAL,
	"sensor_id" INTEGER,
	"value" BYTEA,
	"timestamp" TIMESTAMP,
	PRIMARY KEY("reading_id")
);

CREATE TABLE IF NOT EXISTS "honey_harvest" (
	"harvest_id" SERIAL,
	"hive_id" INTEGER,
	"harvest_date" DATE,
	"quantity" FLOAT,
	"quality_grade" VARCHAR,
	PRIMARY KEY("harvest_id")
);

CREATE TABLE IF NOT EXISTS "observation_log" (
	"log_id" SERIAL,
	"hive_id" INTEGER,
	"observation_date" DATE,
	"description" TEXT,
	"recommendations" TEXT,
	PRIMARY KEY("log_id")
);

CREATE TABLE IF NOT EXISTS "maintenance_plan" (
	"plan_id" SERIAL,
	"apiary_id" INTEGER,
	"planned_date" DATE,
	"work_type" VARCHAR,
	"assigned_to" INTEGER,
	"status" VARCHAR(50) NOT NULL DEFAULT 'pending',
	PRIMARY KEY("plan_id")
);

CREATE TABLE IF NOT EXISTS "incident" (
	"incident_id" SERIAL,
	"hive_id" INTEGER,
	"incident_date" DATE,
	"description" TEXT,
	"severity" VARCHAR,
	"actions_taken" TEXT,
	PRIMARY KEY("incident_id")
);

CREATE TABLE IF NOT EXISTS "production_report" (
	"report_id" SERIAL,
	"apiary_id" INTEGER,
	"start_date" DATE,
	"end_date" DATE,
	"total_honey_produced" FLOAT,
	"total_expenses" FLOAT,
	"curated_by" INTEGER,
	PRIMARY KEY("report_id"),
	UNIQUE("apiary_id", "start_date", "end_date")
);

CREATE TABLE IF NOT EXISTS "weather_data" (
	"weather_id" SERIAL,
	"region_id" INTEGER,
	"date" DATE,
	"temperature" FLOAT,
	"humidity" FLOAT,
	"wind_speed" FLOAT,
	"precipitation" FLOAT,
	PRIMARY KEY("weather_id")
);

CREATE TABLE IF NOT EXISTS "user" (
	"user_id" SERIAL,
	"username" VARCHAR NOT NULL UNIQUE,
	"full_name" VARCHAR NOT NULL,
	"role" ROLE,
	"email" VARCHAR NOT NULL UNIQUE,
	"password" VARCHAR NOT NULL,
	"last_login" TIMESTAMP,
	PRIMARY KEY("user_id")
);

CREATE TABLE IF NOT EXISTS "allowed_region" (
	"id" SERIAL NOT NULL,
	"user_id" INTEGER,
	"region_id" INTEGER,
	PRIMARY KEY("id"),
	UNIQUE ("user_id", "region_id")
);

CREATE TABLE IF NOT EXISTS "region_apiary" (
	"id" SERIAL NOT NULL UNIQUE,
	"apiary_id" INTEGER,
	"region_id" INTEGER,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "worker_group" (
    "group_id" SERIAL PRIMARY KEY,
	"manager_id" INTEGER NOT NULL,
    "group_name" VARCHAR NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
UNIQUE ("manager_id", "group_name")
);

CREATE TABLE IF NOT EXISTS "worker_group_member" (
	"id" SERIAL PRIMARY KEY,
	"group_id" INTEGER NOT NULL,
	"worker_id" INTEGER NOT NULL,
	"joined_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE ("group_id", "worker_id")
);

DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='production_report' AND column_name='curated_by'
    ) THEN
        ALTER TABLE "production_report" ADD COLUMN "curated_by" INTEGER;
    END IF;
END $$;

ALTER TABLE
	"worker_group"
ADD
FOREIGN KEY ("manager_id") REFERENCES "user"("user_id") ON UPDATE NO ACTION ON DELETE CASCADE;
	
ALTER TABLE
	"hive"
ADD
	FOREIGN KEY("apiary_id") REFERENCES "apiary"("apiary_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"bee_community"
ADD
	FOREIGN KEY("hive_id") REFERENCES "hive"("hive_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"veterinary_passport"
ADD
	FOREIGN KEY("bee_community_id") REFERENCES "bee_community"("community_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"veterinary_record"
ADD
	FOREIGN KEY("passport_id") REFERENCES "veterinary_passport"("passport_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"sensor"
ADD
	FOREIGN KEY("hive_id") REFERENCES "hive"("hive_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"sensor_reading"
ADD
	FOREIGN KEY("sensor_id") REFERENCES "sensor"("sensor_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"honey_harvest"
ADD
	FOREIGN KEY("hive_id") REFERENCES "hive"("hive_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"observation_log"
ADD
	FOREIGN KEY("hive_id") REFERENCES "hive"("hive_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"maintenance_plan"
ADD
	FOREIGN KEY("apiary_id") REFERENCES "apiary"("apiary_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"incident"
ADD
	FOREIGN KEY("hive_id") REFERENCES "hive"("hive_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"production_report"
ADD
	FOREIGN KEY("apiary_id") REFERENCES "apiary"("apiary_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"production_report"
ADD
	FOREIGN KEY("curated_by") REFERENCES "user"("user_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"weather_data"
ADD
	FOREIGN KEY("region_id") REFERENCES "region"("region_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"allowed_region"
ADD
	FOREIGN KEY("region_id") REFERENCES "region"("region_id") ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE
	"region_apiary"
ADD
	FOREIGN KEY("apiary_id") REFERENCES "apiary"("apiary_id") ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE
	"region_apiary"
ADD
	FOREIGN KEY("region_id") REFERENCES "region"("region_id") ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE
	"allowed_region"
ADD
	FOREIGN KEY("user_id") REFERENCES "user"("user_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"sensor"
ALTER COLUMN
	"last_reading_time"
SET
	DEFAULT now();

ALTER TABLE
	"sensor_reading"
ALTER COLUMN
	"timestamp"
SET
	DEFAULT now();
ALTER TABLE
	"worker_group_member"
ADD
	FOREIGN KEY ("group_id") REFERENCES "worker_group"("group_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"worker_group_member"
ADD
	FOREIGN KEY ("worker_id") REFERENCES "user"("user_id") ON UPDATE NO ACTION ON DELETE CASCADE;
