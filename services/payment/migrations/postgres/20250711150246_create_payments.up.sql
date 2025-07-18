CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,

    user_id UUID NOT NULL,
    plan_id VARCHAR(100) NOT NULL,

    razorpay_order_id VARCHAR(255) NOT NULL UNIQUE,
    razorpay_payment_id VARCHAR(255),
    razorpay_signature VARCHAR(255),

    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),

    currency VARCHAR(3) NOT NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'completed', 'failed')),

    payment_method VARCHAR(50) NOT NULL DEFAULT 'razorpay'
        CHECK (payment_method IN ('razorpay')),

    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes 
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
