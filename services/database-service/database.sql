CREATE USER 'test'@'%' IDENTIFIED BY 'test@2023';
GRANT ALL PRIVILEGES ON *.* TO 'test'@'%';
CREATE DATABASE ecommerce_microservice;
CREATE USER 'sad'@'%' IDENTIFIED BY 'sad@2023';
GRANT ALL PRIVILEGES ON ecommerce_microservice. * TO 'sad'@'%';
