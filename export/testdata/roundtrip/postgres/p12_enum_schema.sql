CREATE SCHEMA crm;
CREATE TYPE crm.lead_stage AS ENUM ('new', 'contacted', 'qualified', 'won', 'lost');
CREATE TYPE public.mood AS ENUM ('sad', 'ok', 'happy');
CREATE TABLE crm.leads (
    id serial PRIMARY KEY,
    stage crm.lead_stage NOT NULL DEFAULT 'new',
    owner_mood mood,
    score int DEFAULT 0,
    created_at timestamp DEFAULT now()
);
COMMENT ON TABLE crm.leads IS 'Sales leads';
COMMENT ON COLUMN crm.leads.score IS 'Lead score 0-100';
