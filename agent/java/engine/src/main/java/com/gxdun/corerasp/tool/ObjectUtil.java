package com.gxdun.corerasp.tool;

public final class ObjectUtil {
    private ObjectUtil(){}

    public static boolean equals(Object a, Object b) {
        return (a == b) || (a != null && a.equals(b));
    }

}
