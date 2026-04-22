<?php 

namespace App\Middleware;
use App\Core\Middleware;

class AuthApiMiddleware extends Middleware
{
    public function handle()
    {
        if (!isset($_SESSION['authenticated']) || !isset($_SESSION['username'])) {
            return $this->jsonResponse([], "Operação não permitida", 401);
        }
    }
}