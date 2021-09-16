package com.gxdun.corerasp.tool;

import com.gxdun.corerasp.request.AbstractRequest;

import java.net.InetAddress;
import java.net.UnknownHostException;

public final class RequestUtil {
    private RequestUtil(){}

    public static String getIpAddr(AbstractRequest request) {
        String ipAddress = request.getRemoteAddr();
        /*
        String ipAddress = request.getHeader("x-forwarded-for");
        if(ipAddress == null || ipAddress.length() == 0 || "unknown".equalsIgnoreCase(ipAddress)) {
            ipAddress = request.getHeader("Proxy-Client-IP");
        }
        if(ipAddress == null || ipAddress.length() == 0 || "unknown".equalsIgnoreCase(ipAddress)) {
            ipAddress = request.getHeader("WL-Proxy-Client-IP");
        }
        if(ipAddress == null || ipAddress.length() == 0 || "unknown".equalsIgnoreCase(ipAddress)) {
            ipAddress = request.getRemoteAddr();
            if("127.0.0.1".equals(ipAddress)|| "0:0:0:0:0:0:0:1".equals(ipAddress)){
                //根据网卡取本机配置的IP
                try {
                    InetAddress inet = InetAddress.getLocalHost();
                    ipAddress= inet.getHostAddress();
                } catch (UnknownHostException e) {
                }
            }
        }
        */
        //对于通过多个代理的情况，第一个IP为客户端真实IP,多个IP按照','分割
        if(ipAddress!=null && ipAddress.length()>15){
            //"***.***.***.***".length() = 15
            int dotIndex = ipAddress.indexOf(',');
            if(dotIndex > 0){
                ipAddress = ipAddress.substring(0,dotIndex);
            }
        }
        return ipAddress;
    }
}
