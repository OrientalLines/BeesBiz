CREATE OR REPLACE FUNCTION add_constraint_if_not_exists(
    p_table_name TEXT,
    p_constraint_name TEXT,
    p_constraint_sql TEXT
) RETURNS VOID AS $$
BEGIN
    EXECUTE format('ALTER TABLE %I ADD CONSTRAINT %I %s', p_table_name, p_constraint_name, p_constraint_sql);
EXCEPTION
    WHEN duplicate_table THEN
        RAISE NOTICE 'Table constraint %.% already exists', p_table_name, p_constraint_name;
    WHEN duplicate_object THEN
        RAISE NOTICE 'Table constraint %.% already exists', p_table_name, p_constraint_name;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    PERFORM add_constraint_if_not_exists('apiary', 'check_establishment_date', 'CHECK ("establishment_date" <= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('hive', 'check_installation_date', 'CHECK ("installation_date" <= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('bee_community', 'check_queen_age', 'CHECK ("queen_age" >= 0)');
    PERFORM add_constraint_if_not_exists('bee_community', 'check_population_estimate', 'CHECK ("population_estimate" >= 0)');
    PERFORM add_constraint_if_not_exists('veterinary_passport', 'check_issue_date', 'CHECK ("issue_date" <= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('veterinary_passport', 'check_last_inspection_date', 'CHECK ("last_inspection_date" <= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('veterinary_record', 'check_record_date', 'CHECK ("record_date" <= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('honey_harvest', 'check_harvest_date', 'CHECK ("harvest_date" <= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('honey_harvest', 'check_quantity', 'CHECK ("quantity" >= 0)');
    PERFORM add_constraint_if_not_exists('observation_log', 'check_observation_date', 'CHECK ("observation_date" <= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('maintenance_plan', 'check_planned_date', 'CHECK ("planned_date" >= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('incident', 'check_incident_date', 'CHECK ("incident_date" <= CURRENT_DATE)');
    PERFORM add_constraint_if_not_exists('production_report', 'check_date_range', 'CHECK ("start_date" <= "end_date")');
    PERFORM add_constraint_if_not_exists('production_report', 'check_total_honey_produced', 'CHECK ("total_honey_produced" >= 0)');
    PERFORM add_constraint_if_not_exists('weather_data', 'check_weather_date', 'CHECK ("date" <= CURRENT_DATE)');
END $$;
