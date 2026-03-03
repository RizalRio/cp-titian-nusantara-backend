CREATE TABLE collaboration_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_name VARCHAR(255) NOT NULL,
    contact_person VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    collaboration_type VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    proposal_file_url VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending', -- pending, reviewed, accepted, rejected
    is_notified BOOLEAN DEFAULT false,
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_collab_req_status ON collaboration_requests(status);