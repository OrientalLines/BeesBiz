-- Add constraints and checks

\c bee_management


ALTER TABLE Apiary ADD CONSTRAINT check_establishment_date CHECK (establishment_date <= CURRENT_DATE);
ALTER TABLE Hive ADD CONSTRAINT check_installation_date CHECK (installation_date <= CURRENT_DATE);
ALTER TABLE BeeCommunity ADD CONSTRAINT check_queen_age CHECK (queen_age >= 0);
ALTER TABLE BeeCommunity ADD CONSTRAINT check_population_estimate CHECK (population_estimate >= 0);
ALTER TABLE VeterinaryPassport ADD CONSTRAINT check_issue_date CHECK (issue_date <= CURRENT_DATE);
ALTER TABLE VeterinaryPassport ADD CONSTRAINT check_last_inspection_date CHECK (last_inspection_date <= CURRENT_DATE);
ALTER TABLE VeterinaryRecord ADD CONSTRAINT check_record_date CHECK (record_date <= CURRENT_DATE);
ALTER TABLE HoneyHarvest ADD CONSTRAINT check_harvest_date CHECK (harvest_date <= CURRENT_DATE);
ALTER TABLE HoneyHarvest ADD CONSTRAINT check_quantity CHECK (quantity >= 0);
ALTER TABLE ObservationLog ADD CONSTRAINT check_observation_date CHECK (observation_date <= CURRENT_DATE);
ALTER TABLE MaintenancePlan ADD CONSTRAINT check_planned_date CHECK (planned_date >= CURRENT_DATE);
ALTER TABLE Incident ADD CONSTRAINT check_incident_date CHECK (incident_date <= CURRENT_DATE);
ALTER TABLE ProductionReport ADD CONSTRAINT check_date_range CHECK (start_date <= end_date);
ALTER TABLE ProductionReport ADD CONSTRAINT check_total_honey_produced CHECK (total_honey_produced >= 0);
ALTER TABLE WeatherData ADD CONSTRAINT check_weather_date CHECK (date <= CURRENT_DATE);
