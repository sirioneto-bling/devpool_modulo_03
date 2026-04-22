<?php

namespace App\Client;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;
use GuzzleHttp\Exception\ConnectException;

abstract class BaseClient
{
    protected Client $httpClient;
    protected string $baseUrl;
    protected array $defaultHeaders = [];
    protected int $timeout = 10;

    public function __construct()
    {
        $this->httpClient = new Client([
            'base_uri' => $this->baseUrl,
            'timeout' => $this->timeout,
            'headers' => $this->defaultHeaders,
        ]);
    }

    protected function get(string $endpoint, array $queryParams = []): array
    {
        return $this->request('GET', $endpoint, ['query' => $queryParams]);
    }

    protected function post(string $endpoint, array $data = []): array
    {
        return $this->request('POST', $endpoint, ['json' => $data]);
    }

    protected function put(string $endpoint, array $data = []): array
    {
        return $this->request('PUT', $endpoint, ['json' => $data]);
    }

    protected function delete(string $endpoint, array $data = []): array
    {
        return $this->request('DELETE', $endpoint, ['json' => $data]);
    }

    private function request(string $method, string $endpoint, array $options = []): array
    {
        try {
            $response = $this->httpClient->request($method, $endpoint, $options);

            $body = $response->getBody()->getContents();

            return [
                'success' => true,
                'statusCode' => $response->getStatusCode(),
                'data' => json_decode($body, true) ?? $body,
            ];
        } catch (ConnectException $e) {
            return [
                'success' => false,
                'statusCode' => 0,
                'error' => 'Falha na conexão: ' . $e->getMessage(),
            ];
        } catch (RequestException $e) {
            $statusCode = $e->hasResponse() ? $e->getResponse()->getStatusCode() : 0;
            $body = $e->hasResponse() ? $e->getResponse()->getBody()->getContents() : null;

            return [
                'success' => false,
                'statusCode' => $statusCode,
                'error' => json_decode($body, true) ?? $e->getMessage(),
            ];
        }
    }
}
