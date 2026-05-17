# Test data seed

Run this script after applying all migrations, including the catalog seed migration.

```powershell
docker cp .\scripts\seed_test_data.sql fulfillment-env-postgres:/tmp/seed_test_data.sql
docker exec -i fulfillment-env-postgres psql -U postgres -d fulfillment-app -v ON_ERROR_STOP=1 -f /tmp/seed_test_data.sql
```

Test users:

```text
admin@example.com
logist@example.com
worker@example.com
client@example.com
```

Password for all users:

```text
Password123!
```
