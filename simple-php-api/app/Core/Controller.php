<?php

namespace App\Core;

use App\Supports\SupportsCripto\Cripto;
use App\Supports\Traits\HttpRequestResponseTrait;

abstract class Controller
{
    use Cripto;
    use HttpRequestResponseTrait;

    /**
     * Executa um middleware antes da ação do controller.
     *
     * @param string|callable $middleware Nome da classe do middleware ou callable
     * @param string $method Método a ser executado (padrão: 'handle')
     */
    protected function middleware($middleware, string $method = 'handle')
    {
        // Se for uma string (nome da classe), instancia e executa
        if (is_string($middleware) && class_exists($middleware)) {
            $instance = new $middleware();
            if (method_exists($instance, $method)) {
                $instance->$method();
            }
            return;
        }

        // Mantém compatibilidade com callable (forma antiga)
        if (is_callable($middleware)) {
            $middleware();
        }
    }
}

