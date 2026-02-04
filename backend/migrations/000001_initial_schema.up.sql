-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Stock Quotes Table
CREATE TABLE IF NOT EXISTS stock_quotes (
    id BIGSERIAL,
    symbol VARCHAR(10) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    open DECIMAL(10,2) NOT NULL,
    high DECIMAL(10,2) NOT NULL,
    low DECIMAL(10,2) NOT NULL,
    close DECIMAL(10,2) NOT NULL,
    volume BIGINT NOT NULL,
    bid DECIMAL(10,2),
    ask DECIMAL(10,2),
    change DECIMAL(10,2),
    change_pct DECIMAL(5,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

-- Sales Table
CREATE TABLE IF NOT EXISTS sales (
    id BIGSERIAL,
    timestamp TIMESTAMPTZ NOT NULL,
    product_id VARCHAR(50) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    customer_id VARCHAR(50),
    region VARCHAR(100),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) DEFAULT 0,
    revenue DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

-- User Events Table
CREATE TABLE IF NOT EXISTS user_events (
    id BIGSERIAL,
    timestamp TIMESTAMPTZ NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    user_id VARCHAR(50),
    session_id VARCHAR(100),
    page VARCHAR(255),
    device VARCHAR(50),
    browser VARCHAR(50),
    country VARCHAR(100),
    city VARCHAR(100),
    referrer VARCHAR(500),
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

-- Financial Metrics Table
CREATE TABLE IF NOT EXISTS financial_metrics (
    id BIGSERIAL,
    timestamp TIMESTAMPTZ NOT NULL,
    metric_type VARCHAR(50) NOT NULL,
    department VARCHAR(100),
    category VARCHAR(100),
    amount DECIMAL(12,2) NOT NULL,
    budget DECIMAL(12,2),
    variance DECIMAL(12,2),
    variance_pct DECIMAL(5,2),
    period VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

-- Convert to TimescaleDB hypertables
SELECT create_hypertable('stock_quotes', 'timestamp', if_not_exists => TRUE);
SELECT create_hypertable('sales', 'timestamp', if_not_exists => TRUE);
SELECT create_hypertable('user_events', 'timestamp', if_not_exists => TRUE);
SELECT create_hypertable('financial_metrics', 'timestamp', if_not_exists => TRUE);

-- Create indexes for stock_quotes
CREATE INDEX IF NOT EXISTS idx_stock_quotes_symbol ON stock_quotes(symbol);
CREATE INDEX IF NOT EXISTS idx_stock_quotes_timestamp ON stock_quotes(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_stock_quotes_symbol_timestamp ON stock_quotes(symbol, timestamp DESC);

-- Create indexes for sales
CREATE INDEX IF NOT EXISTS idx_sales_category ON sales(category);
CREATE INDEX IF NOT EXISTS idx_sales_region ON sales(region);
CREATE INDEX IF NOT EXISTS idx_sales_product_id ON sales(product_id);
CREATE INDEX IF NOT EXISTS idx_sales_timestamp ON sales(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_sales_category_timestamp ON sales(category, timestamp DESC);

-- Create indexes for user_events
CREATE INDEX IF NOT EXISTS idx_user_events_event_type ON user_events(event_type);
CREATE INDEX IF NOT EXISTS idx_user_events_user_id ON user_events(user_id);
CREATE INDEX IF NOT EXISTS idx_user_events_session_id ON user_events(session_id);
CREATE INDEX IF NOT EXISTS idx_user_events_country ON user_events(country);
CREATE INDEX IF NOT EXISTS idx_user_events_timestamp ON user_events(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_user_events_type_timestamp ON user_events(event_type, timestamp DESC);

-- Create indexes for financial_metrics
CREATE INDEX IF NOT EXISTS idx_financial_metrics_metric_type ON financial_metrics(metric_type);
CREATE INDEX IF NOT EXISTS idx_financial_metrics_department ON financial_metrics(department);
CREATE INDEX IF NOT EXISTS idx_financial_metrics_timestamp ON financial_metrics(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_financial_metrics_type_timestamp ON financial_metrics(metric_type, timestamp DESC);

-- Note: Compression policies removed for MVP
-- Can be added later with: ALTER TABLE <table> SET (timescaledb.compress = true);