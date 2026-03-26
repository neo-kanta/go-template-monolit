-- ============================================================
-- Migration: sys_users + sys_login_log
-- Run this in pgAdmin against your ThaiTA database.
-- ============================================================

-- 1. Users table
CREATE TABLE IF NOT EXISTS sys_users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(50) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
    role        VARCHAR(50) NOT NULL,

    -- Audit
    is_active   BOOLEAN DEFAULT true,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    created_by  VARCHAR(50) DEFAULT 'system',
    updated_by  VARCHAR(50) DEFAULT 'system'
);

-- 2. Login audit log
CREATE TABLE IF NOT EXISTS sys_login_log (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(50) NOT NULL,

    -- Audit
    ip_address  VARCHAR(45),
    user_agent  VARCHAR(255),
    success     BOOLEAN NOT NULL,
    message     VARCHAR(255),
    logged_at   TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Seed POC users
--    Passwords are bcrypt cost-10 hashes of the plaintext shown in comments.
--    These match the hardcoded users in api-gateway/internal/http/handlers/auth.go.
--
--    NOTE: If you prefer, the FND service can auto-seed these on startup.
--    The hashes below were generated with: bcrypt.GenerateFromPassword([]byte(password), 10)

-- admin / admin123
INSERT INTO sys_users (username, password, role, created_by, updated_by) VALUES
  ('admin',    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin',   'system', 'system'),
  ('user',     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user',    'system', 'system'),
  ('manager1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'manager', 'system', 'system'),
  ('manager2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'manager', 'system', 'system'),
  ('maker1',   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'maker',   'system', 'system'),
  ('maker2',   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'maker',   'system', 'system'),
  ('checker1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'checker', 'system', 'system'),
  ('checker2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'checker', 'system', 'system'),
  ('viewer1',  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'viewer',  'system', 'system'),
  ('viewer2',  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'viewer',  'system', 'system'),
  ('auditor1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'auditor', 'system', 'system'),
  ('auditor2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'auditor', 'system', 'system')
ON CONFLICT (username) DO NOTHING;

-- Index for login log queries
CREATE INDEX IF NOT EXISTS idx_sys_login_log_username ON sys_login_log (username);
CREATE INDEX IF NOT EXISTS idx_sys_login_log_logged_at ON sys_login_log (logged_at);
