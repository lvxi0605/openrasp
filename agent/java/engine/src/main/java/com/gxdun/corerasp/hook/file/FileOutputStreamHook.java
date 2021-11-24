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
import com.gxdun.corerasp.tool.FileUtil;
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
 * Created by lxk on 6/8/17.
 * <p>
 * 文件输出流 hook 点
 */
@HookAnnotation
public class FileOutputStreamHook extends AbstractClassHook {
    /**
     * (none-javadoc)
     *
     * @see AbstractClassHook#getType()
     */
    @Override
    public String getType() {
        return "writeFile";
    }

    /**
     * (none-javadoc)
     *
     * @see AbstractClassHook#isClassMatched(String)
     */
    @Override
    public boolean isClassMatched(String className) {
        return "java/io/FileOutputStream".equals(className);
    }

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#hookMethod(CtClass)
     */
    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(FileOutputStreamHook.class, "checkWriteFile", "$1", File.class);
        insertBefore(ctClass.getConstructor("(Ljava/io/File;Z)V"), src);
    }

    /**
     * 写文件hook点
     *
     * @param file
     */
    public static void checkWriteFile(File file) {
        if (file != null) {
            HashMap<String, Object> params = new HashMap<String, Object>();
            params.put("path", file.getPath());
            params.put("realpath", FileUtil.getRealPath(file));
            List<String> stackInfo = StackTrace.getParamStackTraceArray();
            params.put("stack", stackInfo);
            HookHandler.doCheck(CheckParameter.Type.WRITEFILE, params);
        }
    }

}
