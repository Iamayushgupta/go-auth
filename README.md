## Different authentication methods in go

## SQL Table connection
```
my sql -u root -p
```

enter your password

```
create database authDB;
```

```
use database authDB;
```

```
CREATE TABLE users (id INT PRIMARY KEY AUTO_INCREMENT,username VARCHAR(255) UNIQUE NOT NULL,password VARCHAR(255) NOT NULL);
```
