-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_apiary_manager ON "apiary"(manager_id);

CREATE INDEX IF NOT EXISTS idx_hive_apiary ON "hive"(apiary_id);

CREATE INDEX IF NOT EXISTS idx_bee_community_hive ON "bee_community"(hive_id);

CREATE INDEX IF NOT EXISTS idx_veterinary_passport_community ON "veterinary_passport"(bee_community_id);

CREATE INDEX IF NOT EXISTS idx_veterinary_record_passport ON "veterinary_record"(passport_id);

CREATE INDEX IF NOT EXISTS idx_sensor_hive ON "sensor"(hive_id);

CREATE INDEX IF NOT EXISTS idx_sensor_reading_sensor ON "sensor_reading"(sensor_id);

CREATE INDEX IF NOT EXISTS idx_honey_harvest_hive ON "honey_harvest"(hive_id);

CREATE INDEX IF NOT EXISTS idx_observation_log_hive ON "observation_log"(hive_id);

CREATE INDEX IF NOT EXISTS idx_maintenance_plan_apiary ON "maintenance_plan"(apiary_id);

CREATE INDEX IF NOT EXISTS idx_incident_hive ON "incident"(hive_id);

CREATE INDEX IF NOT EXISTS idx_production_report_apiary ON "production_report"(apiary_id);

CREATE INDEX IF NOT EXISTS idx_weather_data_region ON "weather_data"(region_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_allowed_region_user_region
ON "allowed_region" (user_id, region_id);

CREATE INDEX IF NOT EXISTS idx_allowed_region_region ON "allowed_region"(region_id);

CREATE INDEX IF NOT EXISTS idx_region_apiary_apiary ON "region_apiary"(apiary_id);

CREATE INDEX IF NOT EXISTS idx_region_apiary_region ON "region_apiary"(region_id);
