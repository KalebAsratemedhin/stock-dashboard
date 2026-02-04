-- Drop tables in reverse order (to handle any dependencies)
DROP TABLE IF EXISTS financial_metrics CASCADE;
DROP TABLE IF EXISTS user_events CASCADE;
DROP TABLE IF EXISTS sales CASCADE;
DROP TABLE IF EXISTS stock_quotes CASCADE;