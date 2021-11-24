/*
 * Copyright 2021 CORE SHIELD Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.gxdun.corerasp.hook.file;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.tool.Reflection;
import com.gxdun.corerasp.tool.StackTrace;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;
import java.io.File;
import java.io.IOException;
import java.util.HashMap;
import java.util.List;

/**
 * nio files delete hook
 * liergou
 * 2020.8.4
 */
@HookAnnotation
public class NioFilesDeleteHook extends AbstractClassHook {
    /**
     * (none-javadoc)
     *
     * @see AbstractClassHook#getType()
     */
    @Override
    public String getType() {
        return "deleteFile";
    }

    /**
     * (none-javadoc)
     *
     * @see AbstractClassHook#isClassMatched(String)
     */
    @Override
    public boolean isClassMatched(String className) {
        return "java/nio/file/Files".equals(className);
    }

    /**
     * (none-javadoc)
     *
     * @see AbstractClassHook#hookMethod(CtClass)
     */
    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(NioFilesDeleteHook.class, "checkDeleteFile", "$1", Object.class);
        insertBefore(ctClass, "delete", "(Ljava/nio/file/Path;)V", src);
        insertBefore(ctClass, "deleteIfExists", "(Ljava/nio/file/Path;)Z", src);
    }

    /**
     * nio files delete hook
     *
     * @param path 文件路径
     */
    public static void checkDeleteFile(Object path) {
        if (path != null) {
            File file=(File) Reflection.invokeMethod(path, "toFile", new Class[]{});
            HashMap<String, Object> params = new HashMap<String, Object>();
            params.put("path", file.getPath());
            try {
                params.put("realpath", file.getCanonicalPath());
            } catch (IOException e) {
                params.put("realpath", file.getAbsolutePath());
            }

            List<String> stackInfo = StackTrace.getParamStackTraceArray();
            params.put("stack", stackInfo);
            HookHandler.doCheck(CheckParameter.Type.DELETEFILE, params);
        }
    }
}
