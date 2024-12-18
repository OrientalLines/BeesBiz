CREATE DATABASE bee_management;

\ c bee_management CREATE TYPE "role" AS ENUM ('ADMIN', 'WORKER');

CREATE TABLE "region" (
	"region_id" SERIAL,
	"name" VARCHAR NOT NULL,
	"climate_zone" VARCHAR,
	PRIMARY KEY("region_id")
);

CREATE TABLE "apiary" (
	"apiary_id" SERIAL,
	"location" VARCHAR NOT NULL,
	"manager_id" INTEGER,
	"establishment_date" DATE,
	PRIMARY KEY("apiary_id")
);

CREATE TABLE "hive" (
	"hive_id" SERIAL,
	"apiary_id" INTEGER,
	"hive_type" VARCHAR,
	"installation_date" DATE,
	"current_status" VARCHAR,
	PRIMARY KEY("hive_id")
);

CREATE TABLE "bee_community" (
	"community_id" SERIAL,
	"hive_id" INTEGER,
	"queen_age" INTEGER,
	"population_estimate" INTEGER,
	"health_status" VARCHAR,
	PRIMARY KEY("community_id")
);

CREATE TABLE "veterinary_passport" (
	"passport_id" SERIAL,
	"community_id" INTEGER,
	"issue_date" DATE,
	"health_status" VARCHAR,
	"last_inspection_date" DATE,
	PRIMARY KEY("passport_id")
);

CREATE TABLE "veterinary_record" (
	"record_id" SERIAL,
	"passport_id" INTEGER,
	"record_date" DATE,
	"description" TEXT,
	"treatment" TEXT,
	PRIMARY KEY("record_id")
);

CREATE TABLE "sensor" (
	"sensor_id" SERIAL,
	"hive_id" INTEGER,
	"sensor_type" VARCHAR,
	"last_reading" BYTEA,
	"last_reading_time" TIMESTAMP,
	PRIMARY KEY("sensor_id")
);

CREATE TABLE "sensor_reading" (
	"reading_id" SERIAL,
	"sensor_id" INTEGER,
	"value" BYTEA,
	"timestamp" TIMESTAMP,
	PRIMARY KEY("reading_id")
);

CREATE TABLE "honey_harvest" (
	"harvest_id" SERIAL,
	"hive_id" INTEGER,
	"harvest_date" DATE,
	"quantity" FLOAT,
	"quality_grade" VARCHAR,
	PRIMARY KEY("harvest_id")
);

CREATE TABLE "observation_log" (
	"log_id" SERIAL,
	"hive_id" INTEGER,
	"observation_date" DATE,
	"description" TEXT,
	"recommendations" TEXT,
	PRIMARY KEY("log_id")
);

CREATE TABLE "maintenance_plan" (
	"plan_id" SERIAL,
	"apiary_id" INTEGER,
	"planned_date" DATE,
	"work_type" VARCHAR,
	"assigned_to" INTEGER,
	PRIMARY KEY("plan_id")
);

CREATE TABLE "incident" (
	"incident_id" SERIAL,
	"hive_id" INTEGER,
	"incident_date" DATE,
	"description" TEXT,
	"severity" VARCHAR,
	"actions_taken" TEXT,
	PRIMARY KEY("incident_id")
);

CREATE TABLE "production_report" (
	"report_id" SERIAL,
	"apiary_id" INTEGER,
	"start_date" DATE,
	"end_date" DATE,
	"total_honey_produced" FLOAT,
	"total_expenses" FLOAT,
	PRIMARY KEY("report_id")
);

CREATE TABLE "weather_data" (
	"weather_id" SERIAL,
	"region_id" INTEGER,
	"date" DATE,
	"temperature" FLOAT,
	"humidity" FLOAT,
	"wind_speed" FLOAT,
	"precipitation" FLOAT,
	PRIMARY KEY("weather_id")
);

CREATE TABLE "user" (
	"user_id" SERIAL,
	"username" VARCHAR NOT NULL UNIQUE,
	"role" ROLE,
	"email" VARCHAR NOT NULL UNIQUE,
	"last_login" TIMESTAMP,
	PRIMARY KEY("user_id")
);

CREATE TABLE "allowed_region" (
	"id" INTEGER NOT NULL UNIQUE,
	"user_id" INTEGER,
	"region_id" INTEGER,
	PRIMARY KEY("id")
);

CREATE TABLE "region_apiary" (
	"id" INTEGER NOT NULL UNIQUE,
	"apiary_id" INTEGER,
	"region_id" INTEGER,
	PRIMARY KEY("id")
);

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
	FOREIGN KEY("community_id") REFERENCES "bee_community"("community_id") ON UPDATE NO ACTION ON DELETE CASCADE;

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
	"weather_data"
ADD
	FOREIGN KEY("region_id") REFERENCES "region"("region_id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE
	"user"
ADD
	FOREIGN KEY("user_id") REFERENCES "allowed_region"("user_id") ON UPDATE NO ACTION ON DELETE NO ACTION;

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