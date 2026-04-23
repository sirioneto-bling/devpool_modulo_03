<?php

namespace App\Core;

use App\Supports\SupportsCripto\Cripto;
use Exception;

abstract class Model
{
    use Cripto;
    protected $db;

    public function __construct()
    {
        $host = $_ENV['DB_HOST'];
        $dbname = $_ENV['DB_DATABASE'];
        $username = $_ENV['DB_USERNAME'];
        $password = $_ENV['DB_PASSWORD'];
        $port = $_ENV['DB_PORT'] ?? '3306';

        if (!$host || !$dbname || !$username || !$password) {
            throw new Exception("As variáveis de ambiente do banco de dados não estão configuradas corretamente.");
        }

        $opcoes = [
            \PDO::ATTR_ERRMODE => \PDO::ERRMODE_EXCEPTION,
            \PDO::MYSQL_ATTR_INIT_COMMAND => "SET NAMES utf8",
        ];

        $dsn = "mysql:host={$host};dbname={$dbname};port={$port}";

        try {
            $this->db = new \PDO($dsn, $username, $password, $opcoes);
        } catch (\PDOException $e) {
            throw new Exception("Erro ao conectar ao banco de dados: " . $e->getMessage());
        }
    }
}

