<?php
ini_set("display_errors","Off");
//dl("php_sockets.dll");

$HOST = '127.0.0.1';
$PORT = 8888;

// 检测连接是否成功
if ($_SERVER['REQUEST_METHOD']==='GET') 
{
    // var_dump(function_exists("socket_create"));
    $sock = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
    @socket_connect($sock, $HOST, $PORT);
    $msg = 'TO:CONNECT';
    socket_write($sock,$msg);
    $res = socket_read($sock,4096);
    echo $res;
    socket_close($sock);
}

// // 代理获取数据
if ($_SERVER['REQUEST_METHOD']==='POST' && !empty($_POST['DataType'])){
    $sock = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
    @socket_connect($sock, $HOST, $PORT);
    socket_set_option($sock,SOL_SOCKET,SO_RCVTIMEO,array("sec"=>3, "usec"=>0 ));
    socket_set_option($socket,SOL_SOCKET,SO_SNDTIMEO,array("sec"=> 3, "usec"=> 0 ) );
    if ($_POST['DataType'] === 'GetData')
    {
        $msg = 'TO:GET';
        socket_write($sock,$msg);
        $res = socket_read($sock,4096);
        if ($res){
            echo $res;
        }else{
            echo "NO DATA";
        }
    }else if ($_POST['DataType'] === 'PostData' && !empty($_POST['Data'])) 
    {   
        $msg = $_POST['Data'];
        $res = socket_write($sock,$msg);
    }
    socket_close($sock);
}

?>