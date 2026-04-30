CREATE TABLE IF NOT EXISTS public.sensors (
    id TEXT NOT NULL, 
    current_temperature SMALLINT,
    target_temperature SMALLINT,
    PRIMARY KEY (id)
);

-- insert sensors
INSERT INTO 
    public.sensors 
    (id, current_temperature, target_temperature)
VALUES 
    ('sensor1', 1, 1),
    ('sensor2', 1, 1),
    ('sensor3', 1, 1),
    ('sensor4', 1, 1)
ON CONFLICT DO NOTHING;