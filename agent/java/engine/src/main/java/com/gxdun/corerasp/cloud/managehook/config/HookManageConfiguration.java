package com.gxdun.corerasp.cloud.managehook.config;

import java.util.HashMap;
import java.util.Map;
import java.util.Set;

import org.apache.commons.lang3.StringUtils;

import com.gxdun.corerasp.cloud.managehook.command.AbstractCommand;
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.tool.annotation.AnnotationScanner;
import com.gxdun.corerasp.tool.annotation.CommandAnnotation;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;


public class HookManageConfiguration 
{
	/*------------------hookCode start--------------------*/

	
	/**
	 * code,enable
	 */
	private static Map<String,Boolean> hookEnableCache = new HashMap<String,Boolean>();
	
	

	

	
	
	/**
	 * 添加编号相关缓存 现只添加注解中含有code的类
	 */
	public static void addHookCodeInformation(Class<? extends AbstractClassHook> class1)
	{
		HookAnnotation hookAnnotation = class1.getAnnotation(HookAnnotation.class);
		if (hookAnnotation == null)
		{
			return;
		}
		String code = hookAnnotation.code();
		if (StringUtils.isEmpty(code))
		{
			return;
		}

		// 默认开启
		hookEnableCache.put(code, true);
	}
	
	
	
	/**
	 * 判断是否开启
	 * @param code code
	 * @return
	 */
	public static boolean getHookEnable(String code)
	{
		Boolean enable = hookEnableCache.get(code);
		if (enable == null)
		{
			return true;
		}
		return enable.booleanValue();
	}
	
	/**
	 * 关闭hook
	 * @param code hook编码
	 */
	public static void closeHook(String code)
	{
		setHookEnable(code,false);
	}
	
	/**
	 * 启用hook
	 * @param code hook编码
	 */
	public static void openHook(String code)
	{
		setHookEnable(code,true);
	}
	
	/**
	 * 设置hook是否启用
	 * @param code
	 * @param enable
	 */
	public static void setHookEnable(String code, boolean enable)
	{
		if (hookEnableCache.containsKey(code))
		{
			hookEnableCache.put(code, enable);
		}
	}
	
	/*------------------hookCode end--------------------*/
	
	
	
	
	
	/*------------------command start--------------------*/
	private static final String COMMANDPAGE = "com.gxdun.corerasp.cloud.managehook.command";
	
	private static Map<String,Class<? extends AbstractCommand>> cmmandTypeCache = new HashMap<String,Class<? extends AbstractCommand>>();

	/**
	 * 初始命令类型缓存
	 */
	@SuppressWarnings("unchecked")
	public static void initCommandCahe()
	{
		Set<Class> set = AnnotationScanner.getClassWithAnnotation(COMMANDPAGE, CommandAnnotation.class);
		if (set == null || set.isEmpty())
		{
			return;
		}
		
		CommandAnnotation lsAnnotation = null;
		String type = null;
		for (Class<?> class1 : set) 
		{
			if (class1.isInstance(AbstractCommand.class))
			{
				lsAnnotation = class1.getAnnotation(CommandAnnotation.class);
				type = lsAnnotation.type();
				if (StringUtils.isEmpty(type))
				{
					continue;
				}
				cmmandTypeCache.put(type, (Class<? extends AbstractCommand>) class1);
			}
		}
	}

	public static Class<? extends AbstractCommand> getCommandClass(String type)
	{
		return cmmandTypeCache.get(type);
	}
	
	/*------------------command end--------------------*/
	
}
