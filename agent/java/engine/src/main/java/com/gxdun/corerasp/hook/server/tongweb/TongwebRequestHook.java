package com.gxdun.corerasp.hook.server.tongweb;

import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

import com.gxdun.corerasp.hook.server.ServerParamHook;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;

/**
 * @description: 插入Tongweb获取request请求处理hook
 * @author: Baimo
 * @create: 2019/06/11
 */
@HookAnnotation
public class TongwebRequestHook extends ServerParamHook {

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#isClassMatched(String)
     */
    @Override
    public boolean isClassMatched(String className) {
        return "com/tongweb/web/oro/Request".equals(className);
    }

    @Override
    protected void hookMethod(CtClass ctClass, String src) throws NotFoundException, CannotCompileException {
        insertBefore(ctClass, "getParameters", null, src);
    }

}
