<?php

class Core
{
    private $controller;
    private $method;
    private $params = array();

    public function __construct()
    {
        $this->checkUri();
    }

    public function run()
    {
        $currentController = $this->getController();
        $c = new $currentController();
        call_user_func_array(array($c, $this->getMethod()), $this->getParams());
    }

    private function checkUri()
    {
        $url = explode("index.php", $_SERVER["PHP_SELF"]);
        $url = end($url);
        
        if (!empty($url)) {
            $url = explode('/', $url);
            array_shift($url);

            $nameSpace = '';
            foreach ($url as $segment) {
                $nameSpace .= $nameSpace ? '\\' . ucfirst($segment) : ucfirst($segment);
                $controller = $this->checkIfControllerExists($nameSpace);
                array_shift($url);
                if($controller) {
                    $this->controller = $controller;
                    break;
                }
            }

            if(!$this->controller) {
                return $this->controller = $this->getRouteNotFoundController();
            }

            if (isset($url[0])) {
                $method = $url[0];
                if (method_exists($this->controller, $method)) {
                    $this->method = $method;
                    array_shift($url);
                }
            }

            if (isset($url[0])) {
                $this->params = array_filter($url);
            }
        } else {
            $this->controller = $this->getDefaultController();
        }
    }

    private function getController()
    {
        if (class_exists($this->controller)) {
            return $this->controller;
        }
        return $this->getDefaultController();
    }

    private function getMethod()
    {
        if ($this->method && method_exists($this->getController(), $this->method)) {
            return $this->method;
        }

        return DEFAULT_METHOD;
    }

    private function getParams()
    {
        return $this->params;
    }

    private function checkIfControllerExists($nameSpace)
    {
        $controller = NAMESPACE_CONTROLLER . $nameSpace . "Controller";
        return class_exists($controller) ? $controller : null;
    }

    private function getDefaultController()
    {
        return NAMESPACE_CONTROLLER . ucfirst(DEFAULT_CONTROLLER) . "Controller";
    }

    private function getRouteNotFoundController()
    {
        return NAMESPACE_CONTROLLER . ROUTE_NOT_FOUND_CONTROLLER . "Controller";
    }
}
