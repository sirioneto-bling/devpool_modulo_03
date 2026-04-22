<?php

namespace App\Supports\SupportsLog;

trait LogFile
{
    public function setlog($msg, $req = null)
    {
        $arquivo  = fopen("log.txt", "a");
        fwrite($arquivo, $msg . "\n");
        
        if ($req)
            fwrite($arquivo, $req . "\n");

        fclose($arquivo);
    }

    public function clear(){
        return file_put_contents("log.txt", "");
    }
}
