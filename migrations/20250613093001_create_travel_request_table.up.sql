DO $$ BEGIN
    CREATE TYPE travel_request_status AS ENUM ('SOLICITED', 'APPROVED', 'CANCELED');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS travel_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    traveler_name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    destination_name VARCHAR(255) NOT NULL,
    departure_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP,
    status travel_request_status NOT NULL,
    canceled_by UUID,
    approved_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    canceled_at TIMESTAMP,
    approved_at TIMESTAMP
);

ALTER TABLE travel_requests 
ADD CONSTRAINT fk_travel_requests_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE travel_requests 
ADD CONSTRAINT fk_travel_requests_canceled_by 
FOREIGN KEY (canceled_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE travel_requests 
ADD CONSTRAINT fk_travel_requests_approved_by 
FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_travel_requests_user_id ON travel_requests(user_id);
CREATE INDEX idx_travel_requests_status ON travel_requests(status);
CREATE INDEX idx_travel_requests_departure_date ON travel_requests(departure_date);
CREATE INDEX idx_travel_requests_return_date ON travel_requests(return_date);
CREATE INDEX idx_travel_requests_created_at ON travel_requests(created_at);
CREATE INDEX idx_travel_requests_canceled_by ON travel_requests(canceled_by);
CREATE INDEX idx_travel_requests_approved_by ON travel_requests(approved_by);

CREATE TRIGGER update_travel_requests_updated_at 
    BEFORE UPDATE ON travel_requests 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE travel_requests 
ADD CONSTRAINT chk_canceled_at_when_canceled 
CHECK (
    (status = 'CANCELED' AND canceled_at IS NOT NULL AND canceled_by IS NOT NULL) OR
    (status != 'CANCELED' AND canceled_at IS NULL)
);

ALTER TABLE travel_requests 
ADD CONSTRAINT chk_approved_at_when_approved 
CHECK (
    (status = 'APPROVED' AND approved_at IS NOT NULL AND approved_by IS NOT NULL) OR
    (status != 'APPROVED' AND approved_at IS NULL)
);