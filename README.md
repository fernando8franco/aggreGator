# Gator CLI

## Requirements

### Database

You need to have PostgreSQL v15 or later installed.

Create the database "gator":
```bash
psql -U your_user -c "CREATE DATABASE gator;"
```

You can create the tables with the file located in /sql/dumps/dump.sql. Download it and then run the following command:
psql -U your_user -d gator -f /dump.sql
```bash
psql -U your_user -d gator -f /dump.sql
```

Now you should have the database created.

### Config File

You need a file named **.gatorconfig.json** in the user directory with the following content.

```bash
{"db_url":"postgres://your_db_user:your_db_user_password@localhost:5432/gator?sslmode=disable","current_user_name":""}
```

## Install de gator CLI

Install de gator CLI with the following command:
```bash
go install github.com/fernando8franco/gator@latest
```

You can check if everything is working with:
```bash
gator register franco
```