<?php

namespace App\Core;
use App\Supports\Traits\HttpRequestResponseTrait;

abstract class Middleware
{
    use HttpRequestResponseTrait;
}