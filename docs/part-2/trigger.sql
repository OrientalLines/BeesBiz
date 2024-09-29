-- Create triggers

\c bee_management

CREATE OR REPLACE FUNCTION update_last_login()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_login = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_last_login
BEFORE UPDATE ON "User"
FOR EACH ROW
EXECUTE FUNCTION update_last_login();

-- Create a function to grant access to all regions for admin users
CREATE OR REPLACE FUNCTION grant_admin_access()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.role = 'admin' THEN
        INSERT INTO UserRegionAccess (user_id, region_id)
        SELECT NEW.id, r.id
        FROM Region r
        ON CONFLICT (user_id, region_id) DO NOTHING;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically grant access to all regions for admin users
CREATE TRIGGER admin_access_trigger
AFTER INSERT OR UPDATE ON "User"
FOR EACH ROW
EXECUTE FUNCTION grant_admin_access();

CREATE OR REPLACE FUNCTION update_population_estimate()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE BeeCommunity
        SET population_estimate = population_estimate + NEW.quantity
        WHERE id = NEW.bee_community_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE BeeCommunity
        SET population_estimate = population_estimate - OLD.quantity
        WHERE id = OLD.bee_community_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_bee_population
AFTER INSERT OR DELETE ON HoneyHarvest
FOR EACH ROW
EXECUTE FUNCTION update_population_estimate();

-- Trigger to automatically create a veterinary passport for new bee communities
CREATE OR REPLACE FUNCTION create_veterinary_passport()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO VeterinaryPassport (bee_community_id, issue_date, last_inspection_date)
    VALUES (NEW.id, CURRENT_DATE, CURRENT_DATE);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER new_bee_community_passport
AFTER INSERT ON BeeCommunity
FOR EACH ROW
EXECUTE FUNCTION create_veterinary_passport();

-- Trigger to update production report when new honey harvest is added
CREATE OR REPLACE FUNCTION update_production_report()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO ProductionReport (apiary_id, start_date, end_date, total_honey_produced)
    VALUES (
        (SELECT apiary_id FROM Hive WHERE id = NEW.hive_id),
        DATE_TRUNC('month', NEW.harvest_date),
        DATE_TRUNC('month', NEW.harvest_date) + INTERVAL '1 month' - INTERVAL '1 day',
        NEW.quantity
    )
    ON CONFLICT (apiary_id, start_date, end_date)
    DO UPDATE SET total_honey_produced = ProductionReport.total_honey_produced + NEW.quantity;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_honey_production
AFTER INSERT ON HoneyHarvest
FOR EACH ROW
EXECUTE FUNCTION update_production_report();
