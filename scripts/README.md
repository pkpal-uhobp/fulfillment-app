# Test data seed

Run this script after applying all migrations, including catalog seed migration.

```powershell
psql -h localhost -p 5433 -U postgres -d fulfillment-app -f scripts/seed_test_data.sql
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
