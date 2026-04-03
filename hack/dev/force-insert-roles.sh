#!/usr/bin/env bash
# Script to force insert roles into database
# This ensures roles exist even if migrations have been run before

set -euo pipefail

CONTAINER="aion-dev-postgres"
DB_USER="aion"
DB_NAME="aion-api"

echo "🔧 Force inserting roles into database..."

# Delete existing roles (careful!)
docker exec -i "$CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Disable foreign key checks temporarily
SET session_replication_role = 'replica';

-- Delete user_roles that reference roles
DELETE FROM aion_api.user_roles;

-- Delete existing roles
DELETE FROM aion_api.roles;

-- Re-enable foreign key checks
SET session_replication_role = 'origin';

-- Insert roles
INSERT INTO aion_api.roles (name, description, is_active) VALUES
    ('owner', 'System owner with highest privileges', true),
    ('admin', 'Administrator with full system access', true),
    ('user', 'Default user role with basic access', true),
    ('blocked', 'Blocked user with no access', true);

-- Verify
SELECT role_id, name, is_active FROM aion_api.roles ORDER BY role_id;
EOF

echo "✅ Roles inserted successfully!"
