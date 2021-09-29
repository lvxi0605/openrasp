package com.gxdun.corerasp.hook.file;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.messaging.LogTool;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.tool.Reflection;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;
import org.apache.commons.lang3.StringUtils;

import java.io.IOException;
import java.util.Arrays;
import java.util.HashMap;

/**
 * @author lvxi
 * @Description: 文件上传中的文件参数检测
 * @date 2019/9/18 15:13
 */
@HookAnnotation
public class FileUploadItemTomcatHook extends AbstractClassHook {

    private static ThreadLocal<Boolean> enableFileUploadHook = new ThreadLocal<Boolean>() {
        @Override
        protected Boolean initialValue() {
            return true;
        }
    };

    @Override
    public boolean isClassMatched(String className) {
        return "org/apache/tomcat/util/http/fileupload/disk/DiskFileItem".equals(className);
    }

    @Override
    public String getType() {
        return "fileUpload";
    }

    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(FileUploadItemTomcatHook.class,
                "checkFileItemWithBytes", "$0,$_", Object.class, byte[].class);
        String afterSrc = getInvokeStaticSrc(FileUploadItemTomcatHook.class, "enableFileUploadHook", "");
        insertAfter(ctClass, "get", "()[B", src);
        src = getInvokeStaticSrc(FileUploadItemTomcatHook.class,
                "checkFileItemWithStream", "$0", Object.class);

        /**
         * write 和 getInputStream 都会调用 get 方法，因为 hook 点内部也要调用 get 方法，为了避免重复进入 hook 点，
         * 使用 enableFileUploadHook 做个标记，使用 insertAfter 在方法退出的时候复位标记
         */
        insertBefore(ctClass, "write", "(Ljava/io/File;)V", src);
        insertAfter(ctClass, "write", "(Ljava/io/File;)V", afterSrc, true);
        src = getInvokeStaticSrc(FileUploadItemTomcatHook.class,
                "checkFileItemWithStream", "$0", Object.class);

        insertBefore(ctClass, "getInputStream", "()Ljava/io/InputStream;", src);
        insertAfter(ctClass, "getInputStream", "()Ljava/io/InputStream;", afterSrc, true);
    }

    public static void checkFileItemWithBytes(Object item, byte[] content) {
        if (enableFileUploadHook.get()) {
            boolean isFormField = (Boolean) Reflection.invokeMethod(item, "isFormField", new Class[]{});
            if (!isFormField) {
                HashMap<String, Object> params = new HashMap<String, Object>();
                String name = Reflection.invokeStringMethod(item, "getFieldName", new Class[]{});
                params.put("name", name != null ? name : "");
                String filename = Reflection.invokeStringMethod(item, "getName", new Class[]{});
                params.put("filename", filename);
                if (content.length > 4 * 1024) {
                    content = Arrays.copyOf(content, 4 * 1024);
                }
                try {
                    params.put("content", new String(content, getCharSet(item)));
                } catch (Exception e) {
                    params.put("content", new String(content));
                    LogTool.traceHookWarn("failed to get content of multipart the file: " + e.getMessage(), e);
                }
                if (!params.isEmpty()) {
                    HookHandler.doCheck(CheckParameter.Type.FILEUPLOAD, params);
                }
            }
        }
    }

    public static void checkFileItemWithStream(Object item) {
        Reflection.invokeMethod(item, "get", new Class[]{});
        enableFileUploadHook.set(false);
    }

    public static void enableFileUploadHook() {
        enableFileUploadHook.set(true);
    }

    private static String getCharSet(Object fileItem) {
        String charSet = Reflection.invokeStringMethod(fileItem, "getCharSet", new Class[]{});
        if (charSet == null) {
            charSet = HookHandler.requestCache.get().getCharacterEncoding();
        }
        if (!StringUtils.isEmpty(charSet)) {
            return charSet;
        } else {
            return "UTF-8";
        }
    }

}
