CREATE TABLE IF NOT EXISTS sync_sessions (
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
    initial_state JSONB,
    final_state JSONB,
    session_status TEXT
);

CREATE INDEX IF NOT EXISTS idx_sessions_scenario_rule_h
    ON sync_sessions (scenario, learn_rule, h);
CREATE INDEX IF NOT EXISTS idx_sessions_scenario_rule_size
    ON sync_sessions (scenario, learn_rule, network_size);

CREATE INDEX IF NOT EXISTS idx_sessions_scenario_rule
    ON sync_sessions (scenario, learn_rule);

-- CREATE INDEX idx_sessions_firstn_lastk
--     ON sessions (first_n, last_k);
