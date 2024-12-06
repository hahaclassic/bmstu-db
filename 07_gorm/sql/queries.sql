create table temp_users (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    registration_date TIMESTAMP,
    birth_date DATE,
    premium BOOLEAN,
    premium_expiration TIMESTAMP
);

drop table temp_users;
