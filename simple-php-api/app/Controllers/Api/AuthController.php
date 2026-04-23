<?php

namespace App\Controllers\Api;

use App\Core\Controller;

class AuthController extends Controller
{
    public function login()
    {
        $this->validateRequestMethods([HTTP_METHOD_POST]);
        $username = "devpool";
        $password = "asdf000";

        $data = $this->getRequestData();

        $user_input = $data['username'];
        $pass_input = $data['password'];

        if ($user_input == $username && $pass_input == $password) {
            $_SESSION['authenticated'] = true;
            $_SESSION['username'] = $user_input;

            return $this->jsonResponse(
                [
                'email' => 'devpool@mail.com',
                'name' => 'DevPool',
                'token' => '1234567890'
                ],
                'Autenticação efetuada com sucesso',
                200
            );
        }

        return $this->jsonResponse([],'Usuário ou senha inválido', 404);
    }

    public function logout()
    {
        $this->validateRequestMethods([HTTP_METHOD_POST]);
        session_unset();
        session_destroy();

        return $this->jsonResponse([], 'Logout efetuado com sucesso', 200);
    }
}
