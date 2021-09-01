package com.gxdun.corerasp.cloud.managehook.command;

import java.util.Map;

public abstract class AbstractCommand 
{
	protected Map<String,Object> params;
	
	
	public abstract void execute() throws Exception;


	public void setParams(Map<String, Object> params) 
	{
		this.params = params;
	}
	
	
}
