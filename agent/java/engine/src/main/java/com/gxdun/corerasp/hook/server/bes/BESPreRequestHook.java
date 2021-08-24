package com.gxdun.corerasp.hook.server.bes;

import com.gxdun.corerasp.hook.server.ServerPreRequestHook;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

/**
 * @description: BES pre-request hook
 * @author: bes
 * @create: 2020/03/20
 */
@HookAnnotation
public class BESPreRequestHook extends ServerPreRequestHook {

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#isClassMatched(String)
     */
    @Override
    public boolean isClassMatched(String className) {
        return className.endsWith("com/bes/enterprise/webtier/connector/CoyoteAdapter");
    }

    /**
     * (none-javadoc)
     *
     * @see ServerPreRequestHook#hookMethod(CtClass, String)
     */
    @Override
    protected void hookMethod(CtClass ctClass, String src) throws NotFoundException, CannotCompileException {
        insertBefore(ctClass, "service", null, src);
    }

}
