CREATE TABLE users (
                       id serial PRIMARY KEY,
                       first_name VARCHAR ( 50 ) NOT NULL,
                       last_name VARCHAR ( 50 ) NOT NULL,
                       email VARCHAR ( 100 ) NOT NULL,
                       password VARCHAR ( 100 ) NOT NULL,
                       access_level INTEGER DEFAULT 1 NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP
);


CREATE TABLE rooms (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP
);


CREATE TABLE reservations (
                              id SERIAL PRIMARY KEY,
                              first_name VARCHAR(50) NOT NULL,
                              last_name VARCHAR(50) NOT NULL,
                              email  VARCHAR(50) UNIQUE NOT NULL,
                              phone  VARCHAR(60) UNIQUE NOT NULL,
                              start_date DATE NOT NULL ,
                              end_date DATE NOT NULL ,
                              processed INTEGER DEFAULT 0,
                              room_id INTEGER NOT NULL ,
                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              updated_at TIMESTAMP
);


CREATE TABLE restrictions (
                              id SERIAL PRIMARY KEY,
                              start_date DATE NOT NULL ,
                              end_date DATE NOT NULL ,
                              reservation_id INTEGER NOT NULL DEFAULT 0,
                              room_id INTEGER NOT NULL ,
                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              updated_at TIMESTAMP
);

--reservation_id: 0 [a not null value] indicates the restriction inserted by the owner
--reservation_id: greater than 0 indicate the restriction inserted by a guest