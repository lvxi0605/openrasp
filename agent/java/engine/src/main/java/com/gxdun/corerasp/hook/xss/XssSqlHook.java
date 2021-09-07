package com.gxdun.corerasp.hook.xss;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.hook.sql.AbstractSqlHook;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.plugin.info.AttackInfo;
import com.gxdun.corerasp.tool.Reflection;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;
import org.apache.commons.lang3.StringUtils;
import org.apache.log4j.Logger;

import java.io.IOException;
import java.util.*;

@HookAnnotation
public class XssSqlHook extends AbstractSqlHook
{
	private static final Logger LOGGER = Logger.getLogger(XssSqlHook.class.getName());

	private String className;

	/**
	 * 字符串字节数
	 */
	private final static  byte STARTBYTE = "'".getBytes()[0];

	  @Override public boolean isClassMatched(String className)
	  {

//	  if ("com/mysql/jdbc/PreparedStatement".equals(className) || "com/mysql/cj/jdbc/ClientPreparedStatement".equals(className))
//	  {
//		  this.type = SqlType.MYSQL;
//		  this.exceptions = new String[] { "java/sql/SQLException" };
//		  this.className = className;
//		  return true;
//	  }

	  	return false;
	  }

	  @Override public String getType() { return "sql"; }



	public static void checkSql(String sql,byte[][] parameterValues)
	{
		if (StringUtils.isEmpty(sql))
		{
			return;
		}
		if (sql.trim().toUpperCase().startsWith("SELECT"))
		{
			return;
		}

		try
		{

			Map<Integer,String> valuesMap = new LinkedHashMap<Integer,String>();

			int n = -1;
			String value = null;
			for(int i=0; i < parameterValues.length ;i ++)
			{
				if (parameterValues[i].length == 0 || parameterValues[i][0] != STARTBYTE)
				{
					continue;
				}
				value = new String(parameterValues[i]);
				valuesMap.put(i,value);
			}
			if (valuesMap.isEmpty())
			{
				return;
			}
			Map<String,Object> params = new HashMap<String,Object>();
			params.put("values", valuesMap.values().toArray());
			// 获取数据
			//HookHandler.dataThreadHook.set(true);
			HookHandler.doCheck(CheckParameter.Type.XSS_SQL, params);
			byte[][] returnValue = getData();

			if (returnValue == null)
			{
				return;
			}
			if (returnValue.length == valuesMap.size())
			{
				int nf = 0;
				int i = 0;
				for (Map.Entry<Integer,String> entry :valuesMap.entrySet())
				{
					nf = entry.getKey();
					parameterValues[nf] = returnValue[i];
					i ++;
				}
			}
		} catch (Exception e)
		{
			e.printStackTrace();
			return;
		}

	}

	public static void checkSql8(Object delegate) {

	  	if (delegate == null)
		{
			return;
		}

	  	try {

			Object query = Reflection.invokeMethod(delegate,"getQuery",new Class[]{});
			Object sql = Reflection.invokeMethod(query,"getOriginalSql",new Class[]{});
			if (sql == null)
			{
				return;
			}
			if (((String)sql).trim().toUpperCase().startsWith("SELECT"))
			{
				return;
			}
			Object bindings = Reflection.invokeMethod(query,"getQueryBindings",new Class[]{});

			if (bindings == null)
			{
				return;
			}
			Object bindGing = Reflection.invokeMethod(bindings,"getBindValues",new Class[]{});
			if (bindGing == null)
			{
				return;
			}

			List<Object> jr = new ArrayList<Object>();
			List<String> values = new ArrayList<String>();
			String value = null;

			Object[] bindValues = (Object[]) bindGing;
			for (int i = 0; i < bindValues.length; i ++)
			{
				Object bind = bindValues[i];
				if (bind == null)
				{
					return;
				}

				Object mysqlType = Reflection.invokeMethod(bind,"getMysqlType",new Class[]{});
				if (mysqlType == null)
				{
					continue;
				}
				String name = (String) Reflection.invokeMethod(mysqlType,"getName",new Class[]{});
				if ("VARCHAR".equals(name) || "TEXT".equals(name) || "CHAR".equals(name))
				{
					jr.add(bind);

					byte[] zjsx = (byte[]) Reflection.invokeMethod(bind,"getByteValue",new Class[]{});
					value = new String(zjsx);
					values.add(value);
				} else
				{
					continue;
				}
			}

			if (values.isEmpty())
			{
				return;
			}
			Map<String,Object> params = new HashMap<String,Object>();
			params.put("values", values);
			//HookHandler.dataThreadHook.set(true);
			HookHandler.doCheck(CheckParameter.Type.XSS_SQL, params);
			byte[][] returnValue = getData();
			if (returnValue.length == jr.size())
			{
				for(int i=0; i < returnValue.length ;i ++)
				{
					Reflection.invokeMethod(jr.get(i),"setByteValue",new Class[]{byte[].class},returnValue[i]);
				}
			}
		} catch (Exception e)
		{
			e.printStackTrace();
			return;
		}

	}

	  @Override protected void hookMethod(CtClass ctClass) throws IOException,
	  CannotCompileException, NotFoundException
	  {

		 if ("com/mysql/jdbc/PreparedStatement".equals(className))
		 {
			 String[] executeFuncDescs = new String[] {"([[B[Ljava/io/InputStream;[Z[I)Lcom/mysql/jdbc/Buffer;"};
			 String checkSqlSrc = getInvokeStaticSrc(XssSqlHook.class, "checkSql","this.originalSql,$1",String.class,byte[][].class);
			 insertBefore(ctClass, "fillSendPacket", checkSqlSrc, executeFuncDescs);
		 } else if("com/mysql/cj/jdbc/ClientPreparedStatement".equals(className))
		  {
			  String[] executeFuncDescs = new String[] {"()Z"};
			  String checkSqlSrc = getInvokeStaticSrc(XssSqlHook.class, "checkSql8","$0",Object.class);
			  insertBefore(ctClass, "execute", checkSqlSrc, executeFuncDescs);
		  }/* else if(className.equals("com/zaxxer/hikari/pool/HikariProxyPreparedStatement"))
		 {
			 String[] executeFuncDescs = new String[] {"()Z"};
			 String checkSqlSrc = getInvokeStaticSrc(XssSqlHook.class, "checkSql8","this.delegate",Object.class);
			 insertBefore(ctClass, "execute", checkSqlSrc, executeFuncDescs);
		 }*/
	  }


//	  @Override
//		protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
//
//			String[] executeFuncDescs = new String[] {
//					"(ILcom/mysql/jdbc/Buffer;ZZ[Lcom/mysql/jdbc/Field;Z)Lcom/mysql/jdbc/ResultSetInternalMethods;","(ILcom/mysql/cj/protocol/Message;ZZLcom/mysql/cj/protocol/ColumnDefinition;Z)Lcom/mysql/cj/jdbc/result/ResultSetInternalMethods;"};
//
//			 String checkSqlSrc = getInvokeStaticSrc(XssSqlHook.class, "checkSql", "$2", Object.class);
//		  //String checkSqlSrc = getInvokeStaticSrc(XssSqlHook.class, "checkSql", "");
//			checkSqlSrc = checkSqlSrc.replace("try {", "try {System.out.println(\"进入\");");
//			insertBefore(ctClass, "executeInternal", checkSqlSrc, executeFuncDescs);
//
//		}
//
//
//
//
//		public static void checkSql(Object sendPacket)
//		{
//			// 5.x Buffer sendPacket, 8.x Message NativePacketPayload sendPacket
//			if (sendPacket == null)
//			{
//				return;
//			}
//
//
//			Class cl = sendPacket.getClass();
//			String str = null;
//			try {
//				Method getMethod = cl.getMethod("getByteBuffer", new Class[] {});
//				byte[] b = (byte[]) getMethod.invoke(sendPacket);
//				str = new String(b);
//
//				Method setMethod = cl.getMethod("setByteBuffer", new Class[] {byte[].class});
//				str = str.replace("123","陆七八");
//				byte[] bc = str.getBytes();
//				setMethod.invoke(sendPacket,bc);
//				// 5 截5位  8 截一位
//
//				System.out.println("sql" + str);
//			} catch (Exception e) {
//				// TODO Auto-generated catch block
//				e.printStackTrace();
//			}
//			Map<String,Object> params = new HashMap<String,Object>();
//			params.put("sql", str);
//			HookHandler.doCheck(CheckParameter.Type.XSS_SQL, params);
//		}

	private static byte[][] getData()
	{
		try
		{
			Object obj = HookHandler.dataThreadHook.get();
			if (obj instanceof Boolean)
			{
				return null;
			}
			List<AttackInfo> attackInfoList = (List<AttackInfo>) obj;
			if (attackInfoList.isEmpty())
			{
				return null;
			}
			JsonObject data = attackInfoList.get(0).getExtras();
			JsonArray array = data.getAsJsonArray("data");
			if (array == null || array.size() == 0)
			{
				return null;
			}
			byte[][] returnByte = new byte[array.size()][];
			for (int i=0; i< returnByte.length ;i ++)
			{
				returnByte[i] = array.get(i).getAsString().getBytes();
			}
			return returnByte;
		} catch (Exception e)
		{
			return null;
		}finally {
			HookHandler.dataThreadHook.remove();
		}

	}
}
