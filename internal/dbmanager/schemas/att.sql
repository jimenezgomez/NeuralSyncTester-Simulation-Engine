CREATE TABLE IF NOT EXISTS attack_sessions (
    id BIGSERIAL PRIMARY KEY,
    network_size INTEGER,
    first_k INTEGER,
    first_n INTEGER,
    last_k INTEGER,
    last_n INTEGER,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ DEFAULT now(),
    stimulate_iterations INTEGER,
    learn_iterations INTEGER,
    k JSONB,
    n JSONB,
    l INTEGER,
    m INTEGER,
    h INTEGER,
    learn_rule TEXT,
    scenario TEXT,
    attack_type TEXT,
    attacker_count INTEGER,
    session_status TEXT
);

CREATE INDEX idx_attack_sessions_scnario_attack_h
    ON attack_sessions (scenario, attack_type, h);

CREATE INDEX idx_attack_sessions_scenario_attack_size
    ON attack_sessions (scenario, attack_type, network_size);

CREATE INDEX idx_attack_sessions_scenario_attack
    ON attack_sessions (scenario, attack_type);

-- CREATE INDEX idx_attack_sessions_firstn_lastk
--     ON attack_sessions (first_n, last_k);
