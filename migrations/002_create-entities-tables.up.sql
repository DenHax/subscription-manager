-- Create users table
CREATE TABLE IF NOT EXISTS subscriptions.users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL
);

-- Create services table
CREATE TABLE IF NOT EXISTS subscriptions.services (
    service_name VARCHAR(255) PRIMARY KEY
);

-- Create subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions.subscriptions (
    subscription_id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES subscriptions.users(user_id),
    FOREIGN KEY (service_name) REFERENCES subscriptions.services(service_name)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions.subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions.subscriptions(service_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date ON subscriptions.subscriptions(start_date);
CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date ON subscriptions.subscriptions(end_date);