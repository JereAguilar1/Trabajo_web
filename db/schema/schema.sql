CREATE TABLE clientes (
	id_cliente SERIAL NOT NULL,
	nombre VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE desarrolladores (
	id_desarrollador SERIAL NOT NULL,
	nombre VARCHAR(255) NOT NULL,
	apellido VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE proyectos (
	id SERIAL NOT NULL,
	nombre VARCHAR(255) NOT NULL,
	descripcion VARCHAR(255) NOT NULL,
	estado VARCHAR(100) NOT NULL,
	gastos DOUBLE PRECISION NOT NULL,
	presupuesto DOUBLE PRECISION NOT NULL,
	tiempo_estimado DATE NOT NULL,
	id_cliente INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE desarrolladores_en_proyectos (
	id_desarrollador INT NOT NULL,
	id_proyecto INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE clientes
	ADD CONSTRAINT pk_cliente
	PRIMARY KEY (id_cliente);

ALTER TABLE desarrolladores
	ADD CONSTRAINT pk_desarrollador
	PRIMARY KEY (id_desarrollador);

ALTER TABLE proyectos
	ADD CONSTRAINT pk_proyecto
	PRIMARY KEY (id);

ALTER TABLE proyectos
	ADD CONSTRAINT fk_proyecto_cliente
	FOREIGN KEY (id_cliente) REFERENCES clientes(id_cliente);

ALTER TABLE desarrolladores_en_proyectos
	ADD CONSTRAINT pk_desarrolladores_en_proyectos
	PRIMARY KEY (id_desarrollador, id_proyecto);

ALTER TABLE desarrolladores_en_proyectos
	ADD CONSTRAINT fk_desarrolladores_en_proyectos_desarrollador
	FOREIGN KEY (id_desarrollador) REFERENCES desarrolladores(id_desarrollador);

ALTER TABLE desarrolladores_en_proyectos
	ADD CONSTRAINT fk_desarrolladores_en_proyectos_proyecto
	FOREIGN KEY (id_proyecto) REFERENCES proyectos(id);

