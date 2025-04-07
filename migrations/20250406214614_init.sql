-- +goose Up
CREATE TABLE applications (
    id UUID NOT NULL default gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE application_profiles (
    id UUID NOT NULL default gen_random_uuid() PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    graph_id UUID NOT NULL
);

CREATE INDEX idx_application_profiles_application_id ON application_profiles(application_id);

-- +goose Down
DROP TABLE IF EXISTS application_profiles;
DROP TABLE IF EXISTS applications;
