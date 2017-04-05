## Run this command to update all migration files.

Example which replaces all ":" with "-"

```
UPDATE migrations SET migration = replace(migration, ':', '-');
```
