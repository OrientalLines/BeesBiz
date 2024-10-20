-- Add indexes for performance

\c bee_management

CREATE INDEX idx_apiary_region ON Apiary(region_id);
CREATE INDEX idx_hive_apiary ON Hive(apiary_id);
CREATE INDEX idx_bee_community_hive ON BeeCommunity(hive_id);
CREATE INDEX idx_user_region_access_user ON UserRegionAccess(user_id);
CREATE INDEX idx_user_region_access_region ON UserRegionAccess(region_id);
CREATE INDEX idx_veterinary_passport_bee_community ON VeterinaryPassport(bee_community_id);
CREATE INDEX idx_veterinary_record_passport ON VeterinaryRecord(passport_id);
CREATE INDEX idx_honey_harvest_hive ON HoneyHarvest(hive_id);
CREATE INDEX idx_observation_log_hive ON ObservationLog(hive_id);
CREATE INDEX idx_observation_log_user ON ObservationLog(user_id);
CREATE INDEX idx_maintenance_plan_hive ON MaintenancePlan(hive_id);
CREATE INDEX idx_incident_hive ON Incident(hive_id);
CREATE INDEX idx_production_report_apiary ON ProductionReport(apiary_id);
CREATE INDEX idx_weather_data_region ON WeatherData(region_id);
