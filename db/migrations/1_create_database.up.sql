SET
SQL_MODE= "NO_AUTO_VALUE_ON_ZERO";

CREATE TABLE users
(
    id            BIGINT      NOT NULL AUTO_INCREMENT,
    created_at    TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP,
    deleted_at    TIMESTAMP,
    first_name    VARCHAR(45) NOT NULL,
    last_name     VARCHAR(45) NOT NULL,
    gender        VARCHAR(45) NOT NULL,
    relation_type VARCHAR(45) NOT NULL,
    dob           TIMESTAMP   NOT NULL,
    parent_id     VARCHAR(45) DEFAULT NULL,
    external_id   VARCHAR(45) NOT NULL,
    status        TINYINT     DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (external_id),
    FOREIGN KEY (parent_id) REFERENCES users (external_id)
);

CREATE TABLE identities
(
    id             BIGINT      NOT NULL AUTO_INCREMENT,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP,
    deleted_at     TIMESTAMP,
    external_id    VARCHAR(45) NOT NULL,
    user_id        VARCHAR(45),
    identity_type  VARCHAR(45),
    identity_value VARCHAR(45),
    status         TINYINT   DEFAULT 1,
    PRIMARY KEY (id),
    UNIQUE (external_id),
    FOREIGN KEY (user_id) REFERENCES users (external_id)
);

CREATE TABLE addresses
(
    id          BIGINT      NOT NULL AUTO_INCREMENT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP,
    external_id VARCHAR(45) NOT NULL,
    user_id     VARCHAR(45),
    line_1      VARCHAR(45) NOT NULL,
    line_2      VARCHAR(45) NOT NULL,
    city        VARCHAR(45) NOT NULL,
    province    VARCHAR(45) NOT NULL,
    country     VARCHAR(45) NOT NULL,
    latitude    DOUBLE      NOT NULL,
    longitude   DOUBLE      NOT NULL,
    zipcode     VARCHAR(45) NOT NULL,
    status      TINYINT   DEFAULT 1,
    PRIMARY KEY (id),
    UNIQUE (external_id),
    FOREIGN KEY (user_id) REFERENCES users (external_id)
);
