CREATE TABLE users
(
	id INT NOT NULL UNIQUE,
	amount BIGINT CHECK (amount >= 0) NOT NULL 
);

CREATE TABLE invoices
(
	user_id INT REFERENCES users(id) NOT NULL,
	service_id INT NOT NULL,
	order_id BIGINT NOT NULL,
	amount BIGINT CHECK (amount >= 0) NOT NULL,
	PRIMARY KEY(service_id, order_id)
);

CREATE TABLE checks
(
	id BIGSERIAL NOT NULL UNIQUE,
	service_id INT NOT NULL,
	amount BIGINT CHECK (amount >= 0) NOT NULL,
	created_at DATE NOT NULL
);

CREATE TABLE transactions
(
	id BIGSERIAL NOT NULL UNIQUE,
	user_id INT REFERENCES users(id) NOT NULL,
	amount BIGINT NOT NULL,
	description TEXT,
	created_at DATE
);