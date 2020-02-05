<%@page import="java.io.*, java.net.*" trimDirectiveWhitespaces="true"%>
<%
    String HOST = "127.0.0.1";
    int PORT = 8888;
    String method = request.getMethod();
    if (method == "GET"){
        try{
            Socket socket = new Socket(HOST, PORT);
            InputStream inSocket = socket.getInputStream();
            OutputStream outSocket = socket.getOutputStream();
            String msg = "TO:CONNECT";
            outSocket.write(msg.getBytes());
            BufferedReader reader = new BufferedReader(new InputStreamReader(inSocket));
            out.print(reader.readLine());
            reader.close();
            outSocket.close();
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
                BufferedReader reader = new BufferedReader(new InputStreamReader(inSocket));
                try {
                    out.print(reader.readLine());
                }catch(java.io.IOException e) {
                    out.print("NO DATA");
                }
                reader.close();
            }else if (request.getParameter("DataType").equals("PostData") && !request.getParameter("Data").equals(null)){
                String msg = request.getParameter("Data");
                outSocket.write(msg.getBytes());
            }
            outSocket.close();
            socket.close();
        }catch(java.net.ConnectException e){
        }
    }
%>