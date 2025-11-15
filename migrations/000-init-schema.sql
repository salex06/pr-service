CREATE TABLE IF NOT EXISTS teams(
    id BIGSERIAL PRIMARY KEY,
    team_name VARCHAR(128) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(32) NOT NULL,
    team_name VARCHAR(128) REFERENCES teams(team_name) NOT NULL,
    is_active BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS pull_requests(
    id BIGSERIAL PRIMARY KEY,
    pull_request_id VARCHAR(255) UNIQUE NOT NULL,
    pull_request_name TEXT NOT NULL,
    author_id VARCHAR(32) REFERENCES users(user_id) NOT NULL,
    pr_status VARCHAR(16) NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    merged_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS assigned_reviewers(
    user_id VARCHAR(255) REFERENCES users(user_id),
    pull_request_id VARCHAR(255) REFERENCES pull_requests(pull_request_id),
    PRIMARY KEY(user_id, pull_request_id)
);