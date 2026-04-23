<?php

namespace App\Models;

use App\Core\Model;
use PDO;

class Exemplo extends Model
{
    public function getAll()
    {
       $sqp = "SELECT * FROM exemplo";
       $qry = $this->db->query($sqp);
       return $qry->fetchAll(\PDO::FETCH_OBJ);
    }

    public function getExempleById(int $id)
    {
        $sql = "SELECT * FROM exemplo WHERE id = :id";
        $qry = $this->db->prepare($sql);
        $qry->bindValue(":id", $id);
        $qry->execute();

        return $qry->fetch(PDO::FETCH_ASSOC);
    }

    public function create(string $nome)
    {
        $sql = "INSERT INTO exemplo SET nome=:nome";
        $qry = $this->db->prepare($sql);
        $qry->bindValue(":nome", $nome);
        $qry->execute();
        $id =  $this->db->lastInsertId();

        return $this->getExempleById($id);
    }

    public function update(int $id, string $nome)
    {
        $sql = "UPDATE exemplo SET nome=:nome WHERE id=:id";
        $qry = $this->db->prepare($sql);
        $qry->bindValue(":id", $id);
        $qry->bindValue(":nome", $nome);
        $qry->execute();

        return $this->getExempleById($id);
    }

    public function delete(int $id)
    {
        $sql = "DELETE FROM exemplo WHERE id=:id";
        $qry = $this->db->prepare($sql);
        $qry->bindValue(":id", $id);
        $qry->execute();

        return true;
    }
}
