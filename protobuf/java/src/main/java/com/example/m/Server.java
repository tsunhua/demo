package com.example.m;

import com.example.m.pb.AnimalProto;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import org.apache.commons.codec.binary.Hex;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.util.List;

public class Server {
    public static void main(String[] args) throws IOException {
        HttpServer server = HttpServer.create(new InetSocketAddress("localhost", 8080), 0);
        server.createContext("/whoami", new WhoamiHandler());
        server.start();
    }

    static class WhoamiHandler implements HttpHandler{
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            List<String> contentTypes = exchange.getRequestHeaders().get("Content-Type");
            if (contentTypes==null || (contentTypes.size()>0 && !contentTypes.get(0).equals("application/x-protobuf"))){
                exchange.sendResponseHeaders(404,0);
                exchange.close();
                return;
            }

            OutputStream outputStream = exchange.getResponseBody();
            AnimalProto.Animal animal  = AnimalProto.Animal.newBuilder().setId(12).setName("Dokky").build();
            byte[] bs = animal.toByteArray();
            System.out.println(Hex.encodeHexString(bs));
            exchange.sendResponseHeaders(200, bs.length);
            outputStream.write(bs);
            outputStream.flush();
            outputStream.close();
            exchange.close();
        }
    }
}


