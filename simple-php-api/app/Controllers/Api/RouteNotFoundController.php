<?php

namespace App\Controllers\Api;

use App\Core\Controller;
use App\Models\Exemplo;

class RouteNotFoundController extends Controller
{

    public function index()
    {
        return $this->jsonResponse([], "Route not found", 404);
    }
}
