-- 1. Функция для получения общего количества меда, собранного в конкретном улье за определенный период
CREATE OR REPLACE FUNCTION get_total_honey_harvested(
    p_hive_id INTEGER,
    p_start_date DATE,
    p_end_date DATE
) RETURNS DECIMAL AS $$
DECLARE
    total_honey DECIMAL;
BEGIN
    SELECT COALESCE(SUM(quantity::DECIMAL), 0)
    INTO total_honey
    FROM honey_harvest
    WHERE hive_id = p_hive_id
    AND harvest_date BETWEEN p_start_date AND p_end_date;

    RETURN total_honey;
END;
$$ LANGUAGE plpgsql;

-- 2. Процедура для добавления нового наблюдения в журнал
CREATE OR REPLACE PROCEDURE add_observation(
    p_hive_id INTEGER,
    p_observation_date DATE,
    p_description TEXT,
    p_recommendations TEXT
) AS $$
BEGIN
    INSERT INTO observation_log (hive_id, observation_date, description, recommendations)
    VALUES (p_hive_id, p_observation_date, p_description, p_recommendations);
END;
$$ LANGUAGE plpgsql;

-- 3. Функция для получения текущего состояния здоровья пчелиной общины
CREATE OR REPLACE FUNCTION get_community_health_status(p_community_id INTEGER)
RETURNS VARCHAR AS $$
DECLARE
    health_status VARCHAR;
BEGIN
    SELECT vp.health_status
    INTO health_status
    FROM veterinary_passport vp
    WHERE vp.community_id = p_community_id
    ORDER BY vp.last_inspection_date DESC
    LIMIT 1;

    RETURN COALESCE(health_status, 'Unknown');
END;
$$ LANGUAGE plpgsql;

-- 4. Процедура для обновления статуса улья
CREATE OR REPLACE PROCEDURE update_hive_status(
    p_hive_id INTEGER,
    p_new_status VARCHAR
) AS $$
BEGIN
    UPDATE hive
    SET current_status = p_new_status
    WHERE hive_id = p_hive_id;
END;
$$ LANGUAGE plpgsql;

-- 5. Функция для расчета средней температуры в регионе за последние N дней
CREATE OR REPLACE FUNCTION get_avg_temperature(
    p_region_id INTEGER,
    p_days INTEGER
) RETURNS DECIMAL AS $$
DECLARE
    avg_temp DECIMAL;
BEGIN
    SELECT AVG(temperature::DECIMAL)
    INTO avg_temp
    FROM weather_data
    WHERE region_id = p_region_id
    AND date >= CURRENT_DATE - p_days;

    RETURN COALESCE(avg_temp, 0);
END;
$$ LANGUAGE plpgsql;

-- 6. Процедура для назначения работника на план обслуживания
CREATE OR REPLACE PROCEDURE assign_maintenance_plan(
    p_plan_id INTEGER,
    p_user_id INTEGER
) AS $$
BEGIN
    UPDATE maintenance_plan
    SET assigned_to = p_user_id
    WHERE plan_id = p_plan_id;
END;
$$ LANGUAGE plpgsql;

-- 7. Функция для проверки доступа пользователя к региону
CREATE OR REPLACE FUNCTION has_region_access(
    p_user_id INTEGER,
    p_region_id INTEGER
) RETURNS BOOLEAN AS $$
DECLARE
    has_access BOOLEAN;
BEGIN
    SELECT EXISTS (
        SELECT 1
        FROM allowed_region
        WHERE user_id = p_user_id AND region_id = p_region_id
    ) INTO has_access;

    RETURN has_access;
END;
$$ LANGUAGE plpgsql;

-- 8. Процедура для регистрации инцидента
CREATE OR REPLACE PROCEDURE register_incident(
    p_hive_id INTEGER,
    p_incident_date DATE,
    p_description TEXT,
    p_severity VARCHAR
) AS $$
BEGIN
    INSERT INTO incident (hive_id, incident_date, description, severity)
    VALUES (p_hive_id, p_incident_date, p_description, p_severity);
END;
$$ LANGUAGE plpgsql;

-- 9. Функция для получения последних показаний датчика
CREATE OR REPLACE FUNCTION get_latest_sensor_reading(p_hive_id INTEGER, p_sensor_type VARCHAR)
RETURNS TABLE (value BYTEA, timestamp TIMESTAMP) AS $$
BEGIN
    RETURN QUERY
    SELECT sr.value, sr.timestamp
    FROM sensor s
    JOIN sensor_reading sr ON s.sensor_id = sr.sensor_id
    WHERE s.hive_id = p_hive_id AND s.sensor_type = p_sensor_type
    ORDER BY sr.timestamp DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- 10. Процедура для создания отчета о производстве
CREATE OR REPLACE PROCEDURE create_production_report(
    p_apiary_id INTEGER,
    p_start_date DATE,
    p_end_date DATE
) AS $$
DECLARE
    total_honey DECIMAL;
    total_expenses DECIMAL;
BEGIN
    -- Расчет общего количества собранного меда
    SELECT COALESCE(SUM(quantity::DECIMAL), 0)
    INTO total_honey
    FROM honey_harvest hh
    JOIN hive h ON hh.hive_id = h.hive_id
    WHERE h.apiary_id = p_apiary_id
    AND hh.harvest_date BETWEEN p_start_date AND p_end_date;

    -- Здесь должен быть расчет общих расходов (пример)
    total_expenses := 1000.00;

    -- Создание отчета
    INSERT INTO production_report (apiary_id, start_date, end_date, total_honey_produced, total_expenses)
    VALUES (p_apiary_id, p_start_date, p_end_date, total_honey, total_expenses);
END;
$$ LANGUAGE plpgsql;
