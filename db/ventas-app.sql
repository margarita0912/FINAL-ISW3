CREATE DATABASE `ventas-app` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `ventas-app`;


CREATE TABLE usuarios (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    clave VARCHAR(255) NOT NULL,
    rol ENUM('vendedor', 'comprador', 'precio') NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);


CREATE TABLE productos (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    costo DECIMAL(10,2) NOT NULL,
    precio DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);


CREATE TABLE compras (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    usuario_id BIGINT UNSIGNED NOT NULL,
    producto_id BIGINT UNSIGNED NOT NULL,
    cantidad INT NOT NULL,
    costo_unit DECIMAL(10,2) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
    FOREIGN KEY (producto_id) REFERENCES productos(id)
);


CREATE TABLE ventas (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    usuario_id BIGINT UNSIGNED NOT NULL,
    producto_id BIGINT UNSIGNED NOT NULL,
    cantidad INT NOT NULL,
    descuento DECIMAL(5,2) DEFAULT 0,
    precio_final DECIMAL(10,2) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
    FOREIGN KEY (producto_id) REFERENCES productos(id)
);

