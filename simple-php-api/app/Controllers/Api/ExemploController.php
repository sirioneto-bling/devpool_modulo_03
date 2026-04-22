<?php

namespace App\Controllers\Api;

use App\Core\Controller;
use App\Middleware\AuthApiMiddleware;
use App\Models\Exemplo;

class ExemploController extends Controller
{
    protected $repository;

    public function __construct()
    {
        $this->middleware(AuthApiMiddleware::class);
        $this->repository = new Exemplo();
    }

    public function index()
    {
        $this->validateRequestMethods(['GET']);
        return $this->jsonResponse($this->repository->getAll());
    }

    public function show (int $id) {
        $example = $this->repository->getExempleById($id);
        if(!$example) {
            return $this->jsonResponse('Not found', [], 404);
        }

        return $this->jsonResponse($example);
    }

    public function store()
    {
        $this->validateRequestMethods(['POST']);
        $data = $this->getRequestData();
        if(!$data['name']) {
            return $this->jsonResponse(
                [
                    "RequiredParams" => ['name' => 'Name is required']
                ], 
                'error', 400);
        }

        $example = $this->repository->create($data['name'] . ' ' . rand());

        return $this->jsonResponse($example);
    }

    public function update(int $id)
    {
        $this->validateRequestMethods(['PUT']);

        $example = $this->repository->getExempleById($id);
        if(!$example) {
            return $this->jsonResponse('Not found', [], 404);
        }

        $data = $this->getRequestData();
        if(!$data['name']) {
            return $this->jsonResponse(
                [
                    "RequiredParams" => ['name' => 'Name is required']
                ],
                'error', 400);
        }

        $example = $this->repository->update($id, $data['name']);

        return $this->jsonResponse($example);
    }

    public function delete(int $id)
    {
        $this->validateRequestMethods(['DELETE']);

        $example = $this->repository->getExempleById($id);
        if(!$example) {
            return $this->jsonResponse('Not found', [], 404);
        }

        $example = $this->repository->delete($id);

        return $this->jsonResponse($example, "", 204);
    }

    public function exampleRequest($param = null, $param2 = null)
    {
        $this->validateRequestMethods(['POST']);
        $data = $this->getRequestData();
        return $this->jsonResponse('sucesso',
        [
            'params' => [$param, $param2],
            "data" => $data,
        ]);
    }
}
