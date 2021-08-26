package com.gxdun.corerasp.cloud.managehook.command;

import com.gxdun.corerasp.cloud.managehook.config.HookManageConfiguration;
import com.gxdun.corerasp.cloud.managehook.config.HookManageConstants;
import com.gxdun.corerasp.tool.annotation.CommandAnnotation;

@CommandAnnotation(type = "hookEnable")
public class HookEnableCommand extends AbstractCommand
{

	private static final String ENABLE = "enable";
	@Override
	public void execute() {
		
		String code = String.valueOf(params.get(HookManageConstants.HOOKCODE));
		if ("null".equals(code))
		{
			return;
		}
		
		Object enableObj = params.get(ENABLE);
		if (enableObj == null)
		{
			return;
		}
		
		boolean enable = false;
		if (enableObj instanceof String)
		{
			if ("true".equals(enableObj))
			{
				enable = true;
			} else if ("false".equals(enableObj))
			{
				// 不为true 与 false 视为参数错误 不执行命令
				return;
			}
		}
		
		HookManageConfiguration.setHookEnable(code, enable);
		
	}

}
