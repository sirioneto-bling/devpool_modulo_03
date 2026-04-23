<?php

namespace App\Controllers\Api;

use App\Core\Controller;
use App\Client\TaskClient;

class TaskController extends Controller
{
    protected TaskClient $taskClient;

    public function __construct()
    {
        $this->taskClient = new TaskClient();
    }

    /**
     * GET /api/task/index
     */
    public function index()
    {
        $this->validateRequestMethods(['GET']);

        $result = $this->taskClient->listTasks();

        if (!$result['success']) {
            return $this->jsonResponse(
                $result['error'] ?? 'Erro ao listar tasks',
                'error',
                $result['statusCode'] ?: 502
            );
        }

        return $this->jsonResponse($result['data'], 'Tasks listadas com sucesso');
    }

    /**
     * POST /api/task/store
     */
    public function store()
    {
        $this->validateRequestMethods(['POST']);

        $data = $this->getRequestData();

        if (empty($data['title'])) {
            return $this->jsonResponse(
                ['title' => 'O campo "title" é obrigatório'],
                'error',
                400
            );
        }

        $result = $this->taskClient->createTask(
            $data['title'],
            $data['description'] ?? ''
        );

        if (!$result['success']) {
            return $this->jsonResponse(
                $result['error'] ?? 'Erro ao criar task',
                'error',
                $result['statusCode'] ?: 502
            );
        }

        return $this->jsonResponse($result['data'], 'Task criada com sucesso', 201);
    }
}
