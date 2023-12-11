CREATE database hacktoon;
use hacktoon;

CREATE TABLE cities (
	city_name varchar(50) primary key,
	population int
);

CREATE TABLE empresas (
	cnpj int primary key,
	nome_empresa varchar(100),
	senha varchar(100),
	alarms bool
);


CREATE TABLE client (
	client_id int not null AUTO_INCREMENT primary key,
	nome_client varchar(100),
	empresa_client varchar(100),
	email varchar(100),
    power_cap float
);

INSERT INTO empresas(cnpj, nome_empresa, senha) VALUES (777,"Gui", "Gui123");

INSERT INTO cities(city_name, population) VALUES ('Moscow', 12506000);

SELECT * FROM empresas; -- WHERE cnpj = 2;

SELECT * FROM client; -- WHERE empresa_client = "Gui";


---------------

DELIMITER //
CREATE TRIGGER set_alarms_to_false
BEFORE INSERT ON empresas
FOR EACH ROW
BEGIN
  IF NEW.alarms IS NULL THEN
     SET NEW.alarms = false;
  END IF;
END;//
DELIMITER ;

---------------

DROP TABLE empresas;
DROP TABLE client;

