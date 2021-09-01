package com.gxdun.corerasp.tool;

import java.util.Collection;
import java.util.Map;

public final class CollectionUtil {
    private CollectionUtil(){}

    public static boolean isEmpty(Collection collection){
        return collection == null || collection.isEmpty();
    }

    public static boolean isNotEmpty(Collection collection){
        return !CollectionUtil.isEmpty(collection);
    }

    public static boolean isEmpty(Map collection){
        return collection == null || collection.isEmpty();
    }

    public static boolean isNotEmpty(Map collection){
        return !CollectionUtil.isEmpty(collection);
    }


}
