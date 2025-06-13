# Database Setup

1. `brew install dbmate mysql`
2. `export DATABASE_URL=mysql://<username>:<password>@127.0.0.1:3306/<database_name>`
3. `dbmate create` - to create the database with the step-2 if not created before.
4. `dbmate new <migration_file_name>` - it will create a new migration file in ./db/migration/ path
5. `dbmate up` - to apply migration from step-4.
6. `mysql -u root -p my_database < ./db/seeds/books.sql` - this will seed the database.

# FORMAT FOR config.json
create a config.json file in the root directly and populate it with required params
```
{
    "app_name": "BooksApi",
    "server": {
        "port": 8080,
        "read_timeout": 15,
        "write_timeout": 20,
        "idle_timeout": 60
    },
    "basic_auth": {
        "username": "username",
        "password": "password"
    },
    "database": {
        "type": "mysql",
        "host": "127.0.0.1",
        "port": 3306,
        "user": "root",
        "password": "password",
        "db_name": "my_book"
    }
}
```






