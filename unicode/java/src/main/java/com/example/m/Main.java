package com.example.m;

import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

public class Main {
    public static void main(String[] args) {
        printWord("\uD883\uDEDE");
        printWord("𠀾");
        printWord("里");
    }

    private static void printWord(String character) {
        System.out.printf("Character: %s\n", character);
        System.out.printf("Java Style String: %s\n", toJavaString(character));
        System.out.printf("UTF-8 Hex: %s\n", toUtf8Hex(character));
        System.out.printf("UTF-16 Hex: %s\n", toUtf16Hex(character));
        System.out.printf("Unicode String: %s\n", toUnicodeString(character));
    }

    private  static String toJavaString(String character){
        byte[] arr = StandardCharsets.UTF_16BE.encode(character).array();
        switch (arr.length){
            case 2:
                return String.format("\\u%X%X",arr[0],arr[1]);
            case 4:
                return String.format("\\u%X%X\\u%X%X",arr[0],arr[1],arr[2],arr[3]);
            default:
                return "";
        }
    }

    private static String toUtf16Hex(String character){
        byte[] arr = StandardCharsets.UTF_16BE.encode(character).array();
        switch (arr.length){
            case 2:
                return String.format("%X%X",arr[0],arr[1]);
            case 4:
                return String.format("%X%X %X%X",arr[0],arr[1],arr[2],arr[3]);
            default:
                return "";
        }
    }

   private static String toUtf8Hex(String character){
       ByteBuffer buffer = StandardCharsets.UTF_8.encode(character);
       byte[] arr = buffer.array();
        StringBuilder sb = new StringBuilder();
        for (int i=0;i< buffer.remaining();i++){
            sb.append(String.format("%02X", arr[i]));
        }
        return sb.toString();
   }

    private static String toUnicodeString(String character){
        ByteBuffer buffer = Charset.forName("UTF-32").encode(character);
        byte[] arr =buffer.array();
        if (arr[0]==0 && arr[1]==0) {
            return String.format("U+%x%02x",arr[2],arr[3]);
        }else if (arr[0]==0){
            return String.format("U+%x%02x%02x",arr[1],arr[2],arr[3]);
        }else{
            return String.format("U+%x%02x%02x%02x",arr[0],arr[1],arr[2],arr[3]);
        }
    }
}
