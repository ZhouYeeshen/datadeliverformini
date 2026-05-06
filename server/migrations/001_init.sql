-- 优质个体工商企业经营数据上报系统 V2.0
-- 数据库初始化脚本

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 商户表
CREATE TABLE IF NOT EXISTS businesses (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    license_no      VARCHAR(100),
    legal_person    VARCHAR(50) NOT NULL,
    industry_type   VARCHAR(50) NOT NULL,
    address         VARCHAR(500),
    contact_phone   VARCHAR(20),
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_businesses_name ON businesses(name);
CREATE INDEX idx_businesses_industry ON businesses(industry_type);
CREATE INDEX idx_businesses_status ON businesses(status);

-- 微信账号绑定表
CREATE TABLE IF NOT EXISTS wechat_accounts (
    id              BIGSERIAL PRIMARY KEY,
    openid          VARCHAR(100) NOT NULL UNIQUE,
    business_id     BIGINT NOT NULL REFERENCES businesses(id),
    real_name       VARCHAR(50) NOT NULL,
    phone           VARCHAR(20),
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wechat_openid ON wechat_accounts(openid);
CREATE INDEX idx_wechat_business ON wechat_accounts(business_id);

-- 营业数据上报表
CREATE TABLE IF NOT EXISTS reports (
    id                      BIGSERIAL PRIMARY KEY,
    business_id             BIGINT NOT NULL REFERENCES businesses(id),
    report_month            VARCHAR(7) NOT NULL,
    restaurant_revenue      DECIMAL(15,2) NOT NULL DEFAULT 0,
    retail_revenue          DECIMAL(15,2) NOT NULL DEFAULT 0,
    accommodation_revenue   DECIMAL(15,2) NOT NULL DEFAULT 0,
    tobacco_alcohol_revenue DECIMAL(15,2) NOT NULL DEFAULT 0,
    other_revenue           DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_revenue           DECIMAL(15,2) NOT NULL DEFAULT 0,
    source_type             VARCHAR(20) NOT NULL,
    photo_urls              TEXT,
    status                  VARCHAR(20) NOT NULL DEFAULT 'submitted',
    submitted_by            BIGINT DEFAULT 0,
    submitted_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_reports_biz_month ON reports(business_id, report_month) WHERE deleted_at IS NULL;
CREATE INDEX idx_reports_month ON reports(report_month);
CREATE INDEX idx_reports_status ON reports(status);

-- 预警记录表
CREATE TABLE IF NOT EXISTS alerts (
    id              BIGSERIAL PRIMARY KEY,
    business_id     BIGINT NOT NULL REFERENCES businesses(id),
    business_name   VARCHAR(200),
    alert_type      VARCHAR(30) NOT NULL,
    alert_month     VARCHAR(7) NOT NULL,
    reference_value DECIMAL(15,2) DEFAULT 0,
    current_value   DECIMAL(15,2) DEFAULT 0,
    change_percent  DECIMAL(10,2) DEFAULT 0,
    is_read         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_alerts_type ON alerts(alert_type);
CREATE INDEX idx_alerts_business ON alerts(business_id);
CREATE INDEX idx_alerts_read ON alerts(is_read);
CREATE INDEX idx_alerts_month ON alerts(alert_month);

-- 管理员表
CREATE TABLE IF NOT EXISTS admins (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(100) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    role            VARCHAR(30) NOT NULL DEFAULT 'viewer',
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 默认管理员: admin / admin123
-- bcrypt hash for 'admin123'
INSERT INTO admins (username, password_hash, role) VALUES
    ('admin', '$2a$10$rOzR0aQjPMR0qJvL.xm.5uKqjGZ8m5GhqZEpyZqSqOvGZmX6.AUKe', 'admin')
ON CONFLICT (username) DO NOTHING;

-- 收银系统接入凭证表
CREATE TABLE IF NOT EXISTS pos_credentials (
    id              BIGSERIAL PRIMARY KEY,
    business_id     BIGINT NOT NULL UNIQUE REFERENCES businesses(id),
    api_key         VARCHAR(64) NOT NULL UNIQUE,
    secret          VARCHAR(128) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    last_called_at  TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pos_business ON pos_credentials(business_id);
CREATE INDEX idx_pos_apikey ON pos_credentials(api_key);

-- 接口调用日志表
CREATE TABLE IF NOT EXISTS api_call_logs (
    id              BIGSERIAL PRIMARY KEY,
    credential_id   BIGINT REFERENCES pos_credentials(id),
    method          VARCHAR(10),
    path            VARCHAR(200),
    request_body    TEXT,
    response_code   INT,
    ip_address      VARCHAR(50),
    duration_ms     BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_logs_credential ON api_call_logs(credential_id);
CREATE INDEX idx_logs_time ON api_call_logs(created_at DESC);
