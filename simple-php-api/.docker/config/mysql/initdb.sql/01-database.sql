-- ============================================
-- Script de inicialização do banco de dados
-- Este arquivo é executado automaticamente
-- quando o container MySQL é criado pela primeira vez
-- ============================================

CREATE SCHEMA IF NOT EXISTS `devpool_erp` DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci;
USE `devpool_erp`;

-- Tabela de exemplo para testes
CREATE TABLE IF NOT EXISTS `devpool_erp`.`exemplo` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`nome` VARCHAR(255) NOT NULL,
	`dataCriacao` DATETIME NULL DEFAULT CURRENT_TIMESTAMP(),
	PRIMARY KEY (`id`)
);

-- ============================================
-- Dados iniciais para testes
-- Estes registros permitem testar a API imediatamente
-- ============================================

INSERT INTO `exemplo` (`nome`) VALUES 
    ('Primeiro registro de exemplo'),
    ('Segundo registro de exemplo'),
    ('Terceiro registro de exemplo'),
    ('Quarto registro de exemplo'),
    ('Quinto registro de exemplo');

-- ============================================
-- ADICIONE SUAS TABELAS E DADOS ABAIXO
-- ============================================