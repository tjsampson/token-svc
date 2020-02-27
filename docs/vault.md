# Vault

Vault is used to securely store and retrieve secrets, keys, tokens, certs, and eventually PKI.

## Database Secrets Engine

The database secrets engine allows us to create dynamic users/passwords, with short TTLs

```bash
vault write database/config/postgres \
    plugin_name=postgresql-database-plugin \
    allowed_roles="postgres-role" \
    connection_url="postgresql://{{username}}:{{password}}@postgres:5432/" \
    username="postgres" \
    password="postgres"
```

```bash
vault write database/roles/postgres-role \
    db_name=postgres \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
        GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
    default_ttl="1h" \
    max_ttl="24h"
```
