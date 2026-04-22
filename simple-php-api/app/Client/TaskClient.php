<?php

namespace App\Client;

class TaskClient extends BaseClient
{
    protected string $baseUrl = 'http://host.docker.internal:8080';
    protected int $timeout = 10;
    protected array $defaultHeaders = [
        'Accept' => 'application/json',
        'Content-Type' => 'application/json',
    ];

    /**
     * Cria uma nova task na API Go.
     */
    public function createTask(string $title, string $description = ''): array
    {
        $data = ['title' => $title];

        if ($description !== '') {
            $data['description'] = $description;
        }

        return $this->post('/v1/tasks', $data);
    }

    /**
     * Lista todas as tasks da API Go.
     */
    public function listTasks(): array
    {
        return $this->get('/v1/tasks');
    }
}
