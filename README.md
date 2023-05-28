Pagination

- https://jonnylangefeld.com/blog/how-to-write-a-go-api-pagination

DB Date type

- https://stackoverflow.com/questions/409286/should-i-use-the-datetime-or-timestamp-data-type-in-mysql
- https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go
- https://stackoverflow.com/questions/3799193/mysql-data-best-way-to-implement-paging
- https://popsql.com/learn-sql/mysql/how-to-query-date-and-time-in-mysql
  https://github.com/go-sql-driver/mysql/issues/1181

GORM

- https://articles.wesionary.team/implementation-of-uuid-and-binary-16-in-gorm-v2-1c329c352c91
- https://sourcegraph.com/github.com/jinzhu/gorm/-/blob/association_test.go?L601
- https://stackoverflow.com/questions/71333110/a-foreign-key-constraint-fails-when-insert-to-table-with-gorm

Money Type

- https://github.com/shopspring/decimal

SQL
    
```
SELECT GREATEST(100, MAX(id) + 1) FROM table_name INTO @autoinc;
ALTER TABLE table_name AUTO_INCREMENT = @autoinc;
```

```
CREATE SCHEMA `hades_db` DEFAULT CHARACTER SET utf8mb4 ;
```

Docker

- `docker-compose -f stack.yml up`
- `docker-compose -f stack.yml down`

NewRelic

- https://one.newrelic.com/nr1-core/install-newrelic/installation-plan?account=3811459&state=bcb0d7b6-fccd-5d97-c4d2-95c8b92778bc
- https://github.com/newrelic/go-agent/issues/462 LOGS


Project

- build

  `C:\Users\nogue\projs\hades_backend\app> go build -v .`

- docker_build:

  `docker build -t hades_api --progress=plain .`

- docker_run:

    `docker run -p 3333:3333 --name hades_api --network="host" hades_api`


