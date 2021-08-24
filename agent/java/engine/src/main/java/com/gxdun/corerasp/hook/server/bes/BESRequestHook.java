package com.gxdun.corerasp.hook.server.bes;

import com.gxdun.corerasp.hook.server.ServerParamHook;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

/**
 * @description: BES request hook
 * @author: bes
 * @create: 2020/03/20
 */
@HookAnnotation
public class BESRequestHook extends ServerParamHook {

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#isClassMatched(String)
     */
    @Override
    public boolean isClassMatched(String className) {
        return "com/bes/enterprise/webtier/connector/Request".equals(className);
    }

    @Override
    protected void hookMethod(CtClass ctClass, String src) throws NotFoundException, CannotCompileException {
        insertAfter(ctClass, "parseParameters", "()V", src);
    }

}
