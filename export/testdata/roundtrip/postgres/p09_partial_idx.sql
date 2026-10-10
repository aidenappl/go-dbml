CREATE TABLE memberships (
    id bigserial PRIMARY KEY,
    org_id bigint NOT NULL,
    user_id bigint NOT NULL,
    role varchar(20) NOT NULL DEFAULT 'member',
    deleted_at timestamptz
);
CREATE UNIQUE INDEX memberships_org_user_active ON memberships (org_id, user_id) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS memberships_user_idx ON memberships (user_id);
COMMENT ON INDEX memberships_user_idx IS 'lookups by user';
