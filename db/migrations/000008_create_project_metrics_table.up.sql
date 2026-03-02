CREATE TABLE project_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    metric_key VARCHAR(100),
    metric_label VARCHAR(255),
    metric_value NUMERIC,
    metric_unit VARCHAR(50),
    "order" INTEGER DEFAULT 0
);

CREATE INDEX idx_project_metrics_project_id ON project_metrics(project_id);