-- Connect to the database
\c bee_management

-- Function to populate the database with test data
CREATE OR REPLACE FUNCTION populate_test_data()
RETURNS VOID AS $$
BEGIN
    -- Insert test data for Region
    INSERT INTO Region (name, climate_zone) VALUES
    ('North', 'Temperate'),
    ('South', 'Mediterranean'),
    ('East', 'Continental'),
    ('West', 'Oceanic');

    -- Insert test data for Apiary
    INSERT INTO Apiary (location, establishment_date, region_id) VALUES
    ('Forest Edge', '2020-05-15', 1),
    ('Meadow', '2019-07-20', 2),
    ('Mountain Side', '2021-03-10', 3),
    ('Coastal Area', '2018-09-01', 4);

    -- Insert test data for Hive
    INSERT INTO Hive (apiary_id, installation_date) VALUES
    (1, '2020-06-01'),
    (1, '2020-06-02'),
    (2, '2019-08-01'),
    (3, '2021-04-01'),
    (4, '2018-10-01');

    -- Insert test data for BeeCommunity
    INSERT INTO BeeCommunity (hive_id, queen_age, population_estimate) VALUES
    (1, 2, 50000),
    (2, 1, 40000),
    (3, 3, 60000),
    (4, 1, 45000),
    (5, 2, 55000);

    -- Insert test data for User
    INSERT INTO "User" (username, password_hash, role) VALUES
    ('admin', 'hashed_password', 'admin'),
    ('user1', 'hashed_password', 'user'),
    ('user2', 'hashed_password', 'user');

    -- Insert test data for UserRegionAccess
    INSERT INTO UserRegionAccess (user_id, region_id) VALUES
    (2, 1),
    (2, 2),
    (3, 3),
    (3, 4);

    -- Insert test data for HoneyHarvest
    INSERT INTO HoneyHarvest (hive_id, harvest_date, quantity) VALUES
    (1, '2021-08-15', 25.5),
    (2, '2021-08-16', 22.0),
    (3, '2021-08-20', 30.5),
    (4, '2021-09-01', 28.0),
    (5, '2021-09-05', 26.5);

    -- Insert test data for ObservationLog
    INSERT INTO ObservationLog (hive_id, user_id, observation_date, notes) VALUES
    (1, 2, '2021-07-01', 'Bees appear healthy and active'),
    (2, 2, '2021-07-02', 'Queen spotted, laying eggs'),
    (3, 3, '2021-07-10', 'Hive population growing steadily'),
    (4, 3, '2021-07-15', 'Some signs of possible mite infestation'),
    (5, 2, '2021-07-20', 'Honey production looks promising');

    -- Insert test data for MaintenancePlan
    INSERT INTO MaintenancePlan (hive_id, planned_date, description) VALUES
    (1, '2021-10-01', 'Winter preparation'),
    (2, '2021-10-02', 'Winter preparation'),
    (3, '2021-10-05', 'Winter preparation'),
    (4, '2021-10-10', 'Mite treatment'),
    (5, '2021-10-15', 'Winter preparation');

    -- Insert test data for Incident
    INSERT INTO Incident (hive_id, incident_date, description) VALUES
    (4, '2021-07-20', 'Mite infestation detected');

    -- Insert test data for WeatherData
    INSERT INTO WeatherData (region_id, date, temperature, humidity, precipitation) VALUES
    (1, '2021-08-01', 25.5, 60.0, 0.0),
    (2, '2021-08-01', 28.0, 55.0, 0.0),
    (3, '2021-08-01', 23.0, 65.0, 5.0),
    (4, '2021-08-01', 22.0, 70.0, 2.0);

END;
$$ LANGUAGE plpgsql;

-- Execute the function to populate the database
SELECT populate_test_data();

-- Confirm population
SELECT 'Test data inserted successfully!' AS result;
