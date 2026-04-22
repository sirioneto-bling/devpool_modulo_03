<?php

namespace App\Supports\Traits;

use Exception;

trait HttpRequestResponseTrait
{
    protected function validateRequestMethods(array $expectedMethods)
    {
        $requestMethod = $_SERVER['REQUEST_METHOD'];
        if (!in_array($requestMethod, $expectedMethods)) {
            return $this->jsonResponse([], "Método inválido. Esperado um dos: " . implode(", ", $expectedMethods), 405);
        }
    }

    protected function jsonResponse($data = [], $message = "", $httpCode = 200)
    {
        ob_start();

        ob_clean();

        http_response_code($httpCode);

        header('Content-Type: application/json');

        $response = [];

        if(!empty($data)) {
            $response['data'] = $data;
        }

        if($message) {
            $response['message'] = $message;
        }

        echo json_encode($response, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);
        
        exit;
    }

    protected function getRequestData()
    {
        $contentType = $_SERVER['CONTENT_TYPE'] ?? '';
        
        $rawBody = file_get_contents('php://input');
        
        $data = [
            'files' => []
        ];
    
        switch (true) {
            case strpos($contentType, 'application/json') !== false:
                $data = array_merge($data, $this->parseJsonData($rawBody));
                break;
    
            case strpos($contentType, 'application/x-www-form-urlencoded') !== false:
                parse_str($rawBody, $data);
                break;
    
            case strpos($contentType, 'multipart/form-data') !== false:
                $data = array_merge($data, $_POST);
                $data['files'] = $_FILES;
                break;
    
            default:
                $data['raw'] = $rawBody;
        }
    
        return $data;
    }
    
    private function parseJsonData(string $rawBody)
    {
        $data = json_decode($rawBody, true);
        
        if (json_last_error() !== JSON_ERROR_NONE) {
            throw new Exception('Invalid JSON data: ' . json_last_error_msg());
        }
        
        return $data;
    }

}
