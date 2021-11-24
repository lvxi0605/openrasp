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
import com.gxdun.corerasp.config.Config;
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.tool.FileUtil;
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
 *  nio files read and list hook
 *  liergou
 *  2020.7.9
 */
@HookAnnotation
public class NioFilesReadHook extends AbstractClassHook {
    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#getType()
     */
    @Override
    public String getType() {
        return "readFile";
    }

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#isClassMatched(String)
     */
    @Override
    public boolean isClassMatched(String className) {
        return "java/nio/file/Files".equals(className);
    }

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#hookMethod(CtClass)
     */
    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(NioFilesReadHook.class, "checkNioReadFile", "$1", Object.class);
        insertBefore(ctClass, "readAllBytes", "(Ljava/nio/file/Path;)[B", src);
        insertBefore(ctClass, "newInputStream", "(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/io/InputStream;", src);

        //读写channel
        //insertBefore(ctClass, "newByteChannel", "(Ljava/nio/file/Path;Ljava/util/Set;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/channels/SeekableByteChannel;", src);
        //insertBefore(ctClass, "newByteChannel", "(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/nio/channels/SeekableByteChannel;", src);
    }

    /**
     * 文件读取hook点
     *  path转file保持原有file逻辑,增加stackInfo
     * @param path 文件路径
     */
    public static void checkNioReadFile(Object path) {
        boolean checkSwitch = Config.getConfig().getPluginFilter();
        File file= (File) Reflection.invokeMethod(path, "toFile", new Class[]{});
        if (checkSwitch && !file.exists()) {
            return;
        }
        String filepath;
        try {
            filepath = file.getCanonicalPath();
        } catch (Exception e) {
            filepath = file.getAbsolutePath();
        }
        if (filepath.endsWith(".class")) {
            return;
        }
        HashMap<String, Object> params = new HashMap<String, Object>();
        params.put("path", file.getPath());
        List<String> stackInfo = StackTrace.getParamStackTraceArray();
        params.put("stack", stackInfo);
        params.put("realpath", FileUtil.getRealPath(file));
        HookHandler.doCheck(CheckParameter.Type.READFILE, params);
    }
}
