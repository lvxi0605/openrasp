package com.gxdun.corerasp.cloud.managehook.command;

import java.util.Map;
import org.apache.log4j.Logger;

import com.gxdun.corerasp.cloud.managehook.command.factory.CommandFactory;

public class CommandContext 
{
	public static final Logger LOGGER = Logger.getLogger(CommandContext.class.getName());
	private Map<String,Object> params;

	public void setParams(Map<String, Object> params) 
	{
		this.params = params;
	}
	
	
	public void action()
	{
		try 
		{
			AbstractCommand command = CommandFactory.createCommand(params);
			if (command == null)
			{
				return;
			}
			
			command.execute();
		} catch (Exception e) 
		{
			e.printStackTrace();
			LOGGER.error("命令执行出错:" + e.getMessage());
		}
	}
}
