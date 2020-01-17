<?php
ini_set("display_errors","On");
//dl("php_sockets.dll");


// 检测连接是否成功
if ($_SERVER['REQUEST_METHOD']==='GET') {
    var_dump(function_exists("socket_create"));
    $sock = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
    @socket_connect($sock, '172.16.155.1', 8888);
    $msg = 'TO:CONNECT';
    socket_write($sock,$msg);
    $res = socket_read($sock,4096);
    echo $res;
    socket_close($sock);
}

// 代理获取数据
if ($_SERVER['REQUEST_METHOD']==='POST') {
    $sock = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
    @socket_connect($sock, '172.16.155.1', 8888);
    $msg = 'TO:GET';
    socket_write($sock,$msg);
    $res = socket_read($sock,4096);
    echo $res;
    socket_close($sock);
}



?>