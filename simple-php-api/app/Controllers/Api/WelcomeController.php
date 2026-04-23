<?php

namespace App\Controllers\Api;

use App\Core\Controller;

/**
 * Controller de boas-vindas da API.
 * 
 * Este controller é exibido na rota raiz (/) e não requer autenticação.
 * Serve como ponto de entrada para apresentar informações básicas da API.
 */
class WelcomeController extends Controller
{
    public function index()
    {
        return $this->jsonResponse([
            'api' => 'API REST PHP com MVC',
            'versao' => '1.0.0',
            'documentacao' => URL_BASE . '/api',
            'endpoints' => [
                'GET /api/exemplo' => 'Lista todos os registros',
                'GET /api/exemplo/show/{id}' => 'Busca um registro por ID',
                'POST /api/exemplo/store' => 'Cria um novo registro',
                'PUT /api/exemplo/update/{id}' => 'Atualiza um registro',
                'DELETE /api/exemplo/delete/{id}' => 'Remove um registro',
            ]
        ], 'Bem-vindo à API!');
    }
}
