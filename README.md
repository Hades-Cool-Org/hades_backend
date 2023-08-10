**Project Readme**

Welcome to the Hades Backend Project! This readme provides essential information, links, and commands to help you navigate and contribute to the project effectively.

## Table of Contents
- [Pagination](#pagination)
- [DB Date Type](#db-date-type)
- [GORM](#gorm)
- [Money Type](#money-type)
- [SQL](#sql)
- [Docker](#docker)
- [NewRelic](#newrelic)
- [Project](#project)

---

### Pagination

For implementing pagination in your Go API, you can refer to the following resource:
- [How to Write a Go API with Pagination](https://jonnylangefeld.com/blog/how-to-write-a-go-api-pagination)

### DB Date Type

Understanding date and time handling in your MySQL database:
- [Should I Use the DATETIME or TIMESTAMP Data Type in MySQL?](https://stackoverflow.com/questions/409286/should-i-use-the-datetime-or-timestamp-data-type-in-mysql)
- [How to Use Dates and Times in Go](https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go)
- [Best Way to Implement Paging with MySQL Data](https://stackoverflow.com/questions/3799193/mysql-data-best-way-to-implement-paging)
- [Querying Date and Time in MySQL](https://popsql.com/learn-sql/mysql/how-to-query-date-and-time-in-mysql)
- [MySQL Issue: Incorrect DATETIME value](https://github.com/go-sql-driver/mysql/issues/1181)

### GORM

Resources for working with GORM, a popular Go ORM library:
- [Implementation of UUID and Binary 16 in GORM](https://articles.wesionary.team/implementation-of-uuid-and-binary-16-in-gorm-v2-1c329c352c91)
- [GORM Association Testing](https://sourcegraph.com/github.com/jinzhu/gorm/-/blob/association_test.go?L601)
- [Foreign Key Constraint Failure with GORM](https://stackoverflow.com/questions/71333110/a-foreign-key-constraint-fails-when-insert-to-table-with-gorm)

### Money Type

Using the `decimal` package for precise handling of money values in Go:
- [Decimal Package for Handling Money](https://github.com/shopspring/decimal)

### SQL

Helpful SQL queries and schema creation commands:
```sql
SELECT GREATEST(100, MAX(id) + 1) FROM table_name INTO @autoinc;
ALTER TABLE table_name AUTO_INCREMENT = @autoinc;

CREATE SCHEMA `hades_db` DEFAULT CHARACTER SET utf8mb4;
```

### Docker

Commands for managing Docker containers:
- Start containers: `docker-compose -f stack.yml up -d`
- Stop containers: `docker-compose -f stack.yml down -d`

### NewRelic

Integration and installation guides for NewRelic monitoring:
- [NewRelic Installation Plan](https://one.newrelic.com/nr1-core/install-newrelic/installation-plan?account=3811459&state=bcb0d7b6-fccd-5d97-c4d2-95c8b92778bc)
- [NewRelic Go Agent Issue](https://github.com/newrelic/go-agent/issues/462) (Logging)

### Project

Useful build and run commands for the project:
- Build: `go build -v .`
- Docker Build: `docker build -t hades_api --progress=plain .`
- Docker Run:
  ```bash
  docker run -p 3333:3333 --name hades_api -d --network="host" --add-host=host.docker.internal:host-gateway hades_api
  docker exec -it docker-db-1 bash
  mysql -u root -p
  ```
  Adjust the timezone with: `sudo timedatectl set-timezone America/Sao_Paulo`

---

Feel free to explore the resources and commands listed above to enhance your understanding and contribution to the Hades Backend Project. If you have any questions or need further assistance, please don't hesitate to ask. Happy coding!