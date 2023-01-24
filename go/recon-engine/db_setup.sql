CREATE TABLE IF NOT EXISTS programs (
    id SERIAL UNIQUE PRIMARY KEY,
    code VARCHAR(50) UNIQUE,
    name VARCHAR(50),
    platform VARCHAR(20),
    platform_url VARCHAR(255),
    status VARCHAR(10),
    safe_harbor VARCHAR(20),
    managed BOOLEAN,
    category VARCHAR(50),
    recon_date TIMESTAMP
);

CREATE TABLE IF NOT EXISTS targets (
    id SERIAL UNIQUE PRIMARY KEY,
    prog_id INT NOT NULL,
    category VARCHAR(50),
    FOREIGN KEY (prog_id)
        REFERENCES programs (id)
);

CREATE TABLE IF NOT EXISTS subdomains (
    id SERIAL UNIQUE PRIMARY KEY,
    prog_id INT NOT NULL,
    subdomain VARCHAR(64),
    FOREIGN KEY (prog_id)
        REFERENCES programs (id)
);

CREATE TABLE IF NOT EXISTS urls (
    id SERIAL UNIQUE PRIMARY KEY,
    prog_id INT NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    status INT NOT NULL,
    response_length INT NOT NULL,
    recon_date TIMESTAMP,
    FOREIGN KEY (prog_id)
        REFERENCES programs (id)
);

CREATE TABLE IF NOT EXISTS services (
    id SERIAL UNIQUE PRIMARY KEY,
    prog_id INT NOT NULL,
    subdomain VARCHAR(64) NOT NULL,
    ip_addr VARCHAR(15),
    port INT NOT NULL,
    protocol VARCHAR(5),
    name VARCHAR(50),
    product VARCHAR(50),
    version VARCHAR(20),
    additionals VARCHAR(50),
    FOREIGN KEY (prog_id)
        REFERENCES programs (id)
)
