CREATE
    OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at
        = NOW();
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE TABLE Industry
(
    id          SERIAL      NOT NULL PRIMARY KEY,
    title       TEXT       NOT NULL
);

CREATE TABLE Product
(
    id           SERIAL      NOT NULL PRIMARY KEY,
    name         TEXT       NOT NULL,
    description  TEXT       NOT NULL
);

CREATE TABLE Category
(
    id          SERIAL      NOT NULL PRIMARY KEY,
    title       TEXT        NOT NULL,
    industry_id INT         NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (industry_id) REFERENCES Industry (id) ON DELETE CASCADE
);

CREATE TABLE Company
(
    id           SERIAL      NOT NULL PRIMARY KEY,
    email        TEXT        NOT NULL UNIQUE,
    password     TEXT        NOT NULL,
    name         TEXT        NOT NULL,
    legal_name   TEXT        NOT NULL,
    itn          TEXT        DEFAULT '',
    psrn         TEXT        DEFAULT '',
    adress       TEXT        DEFAULT '',
    phone        TEXT        DEFAULT '',
    link         TEXT        DEFAULT '',
    category_id  INT,
    docs         TEXT[],
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (category_id) REFERENCES Category (id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON Company
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();





CREATE TABLE Requests
(
    id          SERIAL      NOT NULL PRIMARY KEY,
    title       TEXT        NOT NULL,
    products    INT,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (products) REFERENCES Company (id) ON DELETE CASCADE
);


CREATE TABLE RequestsProducts
(
    id          SERIAL NOT NULL PRIMARY KEY,
    request_id  INT    NOT NULL,
    product_id  INT    NOT NULL,
    FOREIGN KEY (request_id) REFERENCES Requests (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES Product (id) ON DELETE CASCADE
);

INSERT INTO Industry(title) VALUES('Машиностроение');
INSERT INTO Industry(title) VALUES('MetalProd');
INSERT INTO Industry(title) VALUES('OilProd');
INSERT INTO Category(title, industry_id) VALUES('category1', 1);
INSERT INTO Category(title, industry_id) VALUES('category2', 1);
INSERT INTO Category(title, industry_id) VALUES('category3', 2);
INSERT INTO Category(title, industry_id) VALUES('category4', 3);

INSERT INTO Company(email, password, name, legal_name, itn, psrn, adress, phone, link, category_id) VALUES('test2@mail.ru', 'password123','Ромашка','ООО Ромашка','7727563779','1056749631995','Москва','8(915)9999998','yandex.ru',1);
INSERT INTO Company(email, password, name, legal_name, itn, psrn, adress, phone, link, category_id)  VALUES('test3@mail.ru', 'password123','Пчелка','ООО Пчелка','7727563719','1057749331995','Москва','8(915)9999918','yandex.ru',1);
INSERT INTO Company(email, password, name, legal_name, itn, psrn, adress, phone, link, category_id)  VALUES('test4@mail.ru', 'password123','Пример','ООО Пример','7727563729','1057249631995','Москва','8(915)9999928','yandex.ru',1);
INSERT INTO Company(email, password, name, legal_name, itn, psrn, adress, phone, link, category_id)  VALUES('test5@mail.ru', 'password123','Тест','ООО Тест','7727533779','1057741631995','Москва','8(915)9999938','yandex.ru',2);
INSERT INTO Company(email, password, name, legal_name, itn, psrn, adress, phone, link, category_id)  VALUES('test6@mail.ru', 'password123','Последний','ООО Последний','7727463779','1057749631995','Москва','8(915)9399998','yandex.ru',3);


