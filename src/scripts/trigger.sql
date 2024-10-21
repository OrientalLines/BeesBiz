-- Create function to create triggers if they don't exist
CREATE OR REPLACE FUNCTION create_trigger_if_not_exists(
    p_trigger_name TEXT,
    p_table_name TEXT,
    p_trigger_timing TEXT,
    p_trigger_event TEXT,
    p_trigger_function TEXT
) RETURNS VOID AS $$
BEGIN
    -- Check if the trigger already exists
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.triggers
        WHERE trigger_name = p_trigger_name
          AND event_object_table = p_table_name
    ) THEN
        -- Create the trigger if it doesn't exist
        EXECUTE format(
            'CREATE TRIGGER %I
             %s %s ON %I
             FOR EACH ROW
             EXECUTE FUNCTION %s()',
            p_trigger_name,
            p_trigger_timing,
            p_trigger_event,
            p_table_name,
            p_trigger_function
        );
        RAISE NOTICE 'Trigger % created on table %', p_trigger_name, p_table_name;
    ELSE
        RAISE NOTICE 'Trigger % already exists on table %', p_trigger_name, p_table_name;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE OR REPLACE FUNCTION update_last_login()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_login = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

SELECT create_trigger_if_not_exists(
    'update_user_last_login',
    'user',
    'BEFORE',
    'UPDATE',
    'update_last_login'
);

-- Create a function to grant access to all regions for admin users
CREATE OR REPLACE FUNCTION grant_admin_access()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.role = 'admin' THEN
        INSERT INTO "user_region_access" (user_id, region_id)
        SELECT NEW.id, r.id
        FROM "region" r
        ON CONFLICT (user_id, region_id) DO NOTHING;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

SELECT create_trigger_if_not_exists(
    'admin_access_trigger',
    'user',
    'AFTER',
    'INSERT OR UPDATE',
    'grant_admin_access'
);

CREATE OR REPLACE FUNCTION update_population_estimate()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE "bee_community"
        SET population_estimate = population_estimate + NEW.quantity
        WHERE id = NEW.bee_community_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE "bee_community"
        SET population_estimate = population_estimate - OLD.quantity
        WHERE id = OLD.bee_community_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

SELECT create_trigger_if_not_exists(
    'update_bee_population',
    'honey_harvest',
    'AFTER',
    'INSERT OR DELETE',
    'update_population_estimate'
);

-- Trigger to automatically create a veterinary passport for new bee communities
CREATE OR REPLACE FUNCTION create_veterinary_passport()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO "veterinary_passport" (bee_community_id, issue_date, last_inspection_date)
    VALUES (NEW.id, CURRENT_DATE, CURRENT_DATE);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

SELECT create_trigger_if_not_exists(
    'new_bee_community_passport',
    'bee_community',
    'AFTER',
    'INSERT',
    'create_veterinary_passport'
);

-- Trigger to update production report when new honey harvest is added
CREATE OR REPLACE FUNCTION update_production_report()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO "production_report" (apiary_id, start_date, end_date, total_honey_produced)
    VALUES (
        (SELECT apiary_id FROM "hive" WHERE id = NEW.hive_id),
        DATE_TRUNC('month', NEW.harvest_date),
        DATE_TRUNC('month', NEW.harvest_date) + INTERVAL '1 month' - INTERVAL '1 day',
        NEW.quantity
    )
    ON CONFLICT (apiary_id, start_date, end_date)
    DO UPDATE SET total_honey_produced = "production_report".total_honey_produced + NEW.quantity;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

SELECT create_trigger_if_not_exists(
    'update_honey_production',
    'honey_harvest',
    'AFTER',
    'INSERT',
    'update_production_report'
);
