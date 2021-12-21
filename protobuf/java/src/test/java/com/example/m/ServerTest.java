package com.example.m;

import com.example.m.pb.AnimalProto;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.io.IOUtils;
import org.apache.hc.client5.http.classic.methods.HttpGet;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.CloseableHttpResponse;
import org.apache.hc.client5.http.impl.classic.HttpClients;
import org.apache.hc.core5.http.io.entity.EntityUtils;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.io.IOException;

import javax.annotation.processing.SupportedSourceVersion;
import javax.lang.model.SourceVersion;

class ServerTest {

    @Test
    public  void testServer() throws IOException {
        CloseableHttpClient httpClient = HttpClients.createDefault();
        HttpGet httpGet = new HttpGet("http://localhost:8080/whoami");
        httpGet.setHeader("Content-Type","application/x-protobuf");
        CloseableHttpResponse response = httpClient.execute(httpGet);

        byte[] bs = IOUtils.toByteArray(response.getEntity().getContent());
        assertEquals(200, response.getCode());
        assertEquals("080c1205446f6b6b79", Hex.encodeHexString(bs));
        AnimalProto.Animal animal = AnimalProto.Animal.parseFrom(bs);
        assertEquals(12, animal.getId());
        assertEquals("Dokky", animal.getName());
        EntityUtils.consume(response.getEntity());
        httpClient.close();
    }
}
