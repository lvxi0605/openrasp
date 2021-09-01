package com.gxdun.corerasp.cloud.managehook.command.factory;

import java.util.Map;

import org.apache.log4j.Logger;

import com.gxdun.corerasp.cloud.managehook.command.AbstractCommand;
import com.gxdun.corerasp.cloud.managehook.config.HookManageConfiguration;
import com.gxdun.corerasp.cloud.managehook.config.HookManageConstants;

/**
 * 命令工厂
 * @author dell
 *
 */
public class CommandFactory 
{
	public static final Logger LOGGER = Logger.getLogger(CommandFactory.class.getName());
	private CommandFactory()
	{
		
	}
	/**
	 * 创建命令
	 * @param params
	 * @return
	 */
	public static AbstractCommand createCommand(Map<String,Object> params)
	{
		Object typObj = params.get(HookManageConstants.TYPE);
		if (typObj == null)
		{
			return null;
		}
	
		String type = String.valueOf(typObj);
		Class<? extends AbstractCommand> classa = HookManageConfiguration.getCommandClass(String.valueOf(typObj));
		if (classa == null)
		{
			return null;
		}
		AbstractCommand command = null;
		try 
		{
			 command = classa.newInstance();
			 command.setParams(params);
		} catch (Exception e) 
		{
			LOGGER.info("创建命令实体失败 type:" + type);
			return null;
		}
		return command;
	}
	
}
