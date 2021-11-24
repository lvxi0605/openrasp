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

package com.gxdun.corerasp.hook.sql;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.config.Config;
import com.gxdun.corerasp.exceptions.SecurityException;
import com.gxdun.corerasp.messaging.ErrorType;
import com.gxdun.corerasp.messaging.LogTool;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.plugin.checker.local.ConfigurableChecker;
import com.gxdun.corerasp.plugin.info.AttackInfo;
import com.gxdun.corerasp.plugin.info.EventInfo;
import com.gxdun.corerasp.tool.CollectionUtil;
import com.gxdun.corerasp.tool.ObjectUtil;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import com.gxdun.corerasp.hook.model.SQLPreparedParam;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

import java.io.IOException;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.util.*;

/**
 * Created by tyy on 18-4-28.
 * <P>
 * sql Prepare 查询 hook 点
 */
@HookAnnotation
public class SQLPreparedStatementHook extends AbstractSqlHook {

    private String className;

    public static ThreadLocal<WeakHashMap<PreparedStatement, Map<Integer, SQLPreparedParam>>> sqlParamsData = new ThreadLocal<WeakHashMap<PreparedStatement, Map<Integer, SQLPreparedParam>>>();

    @Override
    public boolean isClassMatched(String className) {

        /* MySQL */
        if ("com/mysql/jdbc/PreparedStatement".equals(className)
                || "com/mysql/cj/jdbc/PreparedStatement".equals(className)
                || "com/mysql/cj/jdbc/ClientPreparedStatement".equals(className)) {
            this.type = SqlType.MYSQL;
            this.exceptions = new String[]{"java/sql/SQLException"};
            // 钩子报错，暂时关闭
            return true;
        }

        /* SQLite */
        if ("org/sqlite/PrepStmt".equals(className)
                || "org/sqlite/jdbc3/JDBC3PreparedStatement".equals(className)) {
            this.type = SqlType.SQLITE;
            this.exceptions = new String[]{"java/sql/SQLException"};
            return true;
        }

        /* Oracle */
        if ("oracle/jdbc/driver/OraclePreparedStatement".equals(className)) {
            this.type = SqlType.ORACLE;
            this.exceptions = new String[]{"java/sql/SQLException"};
            return true;
        }

        /* SQL Server */
        if ("com/microsoft/sqlserver/jdbc/SQLServerPreparedStatement".equals(className)) {
            this.type = SqlType.SQLSERVER;
            this.exceptions = new String[]{"com/microsoft/sqlserver/jdbc/SQLServerException"};
            return true;
        }

        /* PostgreSQL */
        if ("org/postgresql/jdbc/PgPreparedStatement".equals(className)
                || "org/postgresql/jdbc1/AbstractJdbc1Statement".equals(className)
                || "org/postgresql/jdbc2/AbstractJdbc2Statement".equals(className)
                || "org/postgresql/jdbc3/AbstractJdbc3Statement".equals(className)
                || "org/postgresql/jdbc3g/AbstractJdbc3gStatement".equals(className)
                || "org/postgresql/jdbc4/AbstractJdbc4Statement".equals(className)) {
            this.className = className;
            this.type = SqlType.PGSQL;
            this.exceptions = new String[]{"java/sql/SQLException"};
            return true;
        }

        /* HSqlDB */
        if ("org/hsqldb/jdbc/JDBCPreparedStatement".equals(className)
                || "org/hsqldb/jdbc/jdbcPreparedStatement".equals(className)) {
            this.type = SqlType.HSQL;
            this.exceptions = new String[]{"java/sql/SQLException"};
            return true;
        }

        return false;
    }

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#getType()
     */
    @Override
    public String getType() {
        return "sqlPrepared";
    }

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#hookMethod(CtClass)
     */
    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        hookSqlPreparedStatementMethod(ctClass);
    }

    private void hookSqlPreparedStatementMethod(CtClass ctClass) throws NotFoundException, CannotCompileException {
        String originalSqlCode = null;
//        String checkSqlSrc = null;
        if (SqlType.MYSQL.equals(this.type)) {
            if ("com/mysql/cj/jdbc/ClientPreparedStatement".equals(className)) {
                originalSqlCode = "((com.mysql.cj.PreparedQuery)this.query).getOriginalSql()";
            } else {
                originalSqlCode = "originalSql";
            }
        } else if (SqlType.SQLITE.equals(this.type)
                || SqlType.HSQL.equals(this.type)) {
            originalSqlCode = "this.sql";
        } else if (SqlType.SQLSERVER.equals(this.type)) {
            originalSqlCode = "preparedSQL";
        } else if (SqlType.PGSQL.equals(this.type)) {
            if ("org/postgresql/jdbc/PgPreparedStatement".equals(className)) {
                originalSqlCode = "preparedQuery.query.toString(preparedQuery.query.createParameterList())";
            } else {
                originalSqlCode = "preparedQuery.toString(preparedQuery.createParameterList())";
            }
        } else if (SqlType.ORACLE.equals(this.type)) {
            originalSqlCode = "this.sqlObject.getOriginalSql()";
        }
        if (originalSqlCode != null) {
//            checkSqlSrc = getInvokeStaticSrc(SQLStatementHook.class, "checkSQL",
//                    "\"" + type.name + "\"" + ",$0," + originalSqlCode, String.class, Object.class, String.class);
//            insertBefore(ctClass, "execute", "()Z", checkSqlSrc);
//            insertBefore(ctClass, "executeUpdate", "()I", checkSqlSrc);
//            insertBefore(ctClass, "executeQuery", "()Ljava/sql/ResultSet;", checkSqlSrc);
//            try {
//                insertBefore(ctClass, "executeBatch", "()[I", checkSqlSrc);
//            } catch (CannotCompileException e) {
//                insertBefore(ctClass, "executeBatchInternal", null, checkSqlSrc);
//            }
            insertSetSqlPreparedStatementParams(ctClass,"setString","(ILjava/lang/String;)V",originalSqlCode);
            insertSetSqlPreparedStatementParams(ctClass,"setNString","(ILjava/lang/String;)V",originalSqlCode);
            addCheckSQLXSS(ctClass,originalSqlCode);

            addCatch(ctClass, "execute", null, originalSqlCode);
            addCatch(ctClass, "executeUpdate", null, originalSqlCode);
            addCatch(ctClass, "executeQuery", null, originalSqlCode);
            try {
                addCatch(ctClass, "executeBatch", null, originalSqlCode);
            } catch (CannotCompileException e) {
                addCatch(ctClass, "executeBatchInternal", null, originalSqlCode);
            }
        }
//        else if (SQL_TYPE_DB2.equals(this.type)) {
//            checkSqlSrc = getInvokeStaticSrc(SQLStatementHook.class, "checkSQL",
//                    "\"" + type.name + "\"" + ",$0,$1", String.class, Object.class, String.class);
//            insertBefore(ctClass, "prepareStatement", null, checkSqlSrc);
//        }
    }

    private void addCheckSQLXSS(CtClass ctClass,String originalSqlCode)  throws NotFoundException, CannotCompileException {

        String scr = getInvokeStaticSrc(SQLPreparedStatementHook.class, "checkSQLXSS","\"" + type.name + "\",$0,"+originalSqlCode, PreparedStatement.class, int.class, String.class);

        insertBefore(ctClass, "execute", null, scr);
        insertBefore(ctClass, "executeUpdate", null, scr);
        insertBefore(ctClass, "executeQuery", null, scr);
        try {
            insertBefore(ctClass, "executeBatch", null, scr);
        } catch (CannotCompileException e) {
            insertBefore(ctClass, "executeBatchInternal", null, scr);
        }
    }

    public static void checkSQLXSS(String server, PreparedStatement preparedStatement, String stmt) {
        if (stmt != null && !stmt.isEmpty() ) {
            Map<Integer, SQLPreparedParam>  preparedStatementSqlParams = SQLPreparedStatementHook.getAndRemovePreparedStatementSqlParams(preparedStatement);
            if(CollectionUtil.isEmpty(preparedStatementSqlParams)){
                return;
            }
            if(!SQLPreparedStatementHook.isChangeDbSQL(stmt)){
                return;
            }
            HashMap<String, Object> params = new HashMap<String, Object>(4);
            params.put("server", server);
            params.put("query", stmt);
            params.put("preparedStatementSqlParams",preparedStatementSqlParams);
            HookHandler.doCheck(CheckParameter.Type.XSS_SQL, params);
            List<EventInfo> eventInfos = HookHandler.dataThreadHook.get();
            if(CollectionUtil.isNotEmpty(eventInfos)){
                for (EventInfo eventInfo : eventInfos) {
                   if(!(eventInfo instanceof AttackInfo)) {
                       continue;
                   }
                   AttackInfo attackInfo = (AttackInfo) eventInfo;
                   JsonObject extras = attackInfo.getExtras();
                   if(extras==null){
                       continue;
                   }
                   JsonArray xssSqlChangeParams = extras.getAsJsonArray("xssSqlChangeParams");
                   if(xssSqlChangeParams==null ||  xssSqlChangeParams.size()==0){
                       continue;
                   }
                    for (JsonElement xssSqlChangeParam : xssSqlChangeParams) {
                        JsonObject paramJsonObject = xssSqlChangeParam.getAsJsonObject();
                        int index = paramJsonObject.get("index").getAsInt();
                        String method = paramJsonObject.get("method").getAsString();
                        String value = paramJsonObject.get("value").getAsString();
                        if("setString".equals(method)){
                            try {
                                preparedStatement.setString(index,value);
                            } catch (SQLException e) {
                                LogTool.error(ErrorType.PLUGIN_ERROR, "xss_sql setString 错误", e);
                            }
                        }else if ("setNString".equals(method)){
                            try {
                                preparedStatement.setNString(index,value);
                            } catch (SQLException e) {
                                LogTool.error(ErrorType.PLUGIN_ERROR, "xss_sql setNString 错误", e);
                            }
                        }
                    }
                }
            }
        }
    }

    private static boolean isChangeDbSQL(String stmt){
        String tmpStmt = stmt.trim();
        if(tmpStmt.length()<6){
            return false;
        }
        String sqlType = tmpStmt.substring(0, 6);
        if("select".equalsIgnoreCase(sqlType) || "delete".equalsIgnoreCase(sqlType)){
            return false;
        }
        return true;
    }


    private void insertSetSqlPreparedStatementParams(CtClass ctClass,String methodName, String desc,String originalSqlCode){
        try {
            String putScr = getInvokeStaticSrc(SQLPreparedStatementHook.class, "putPreparedStatementSqlParam","$0,"+originalSqlCode+",\""+methodName+"\",$1,$2", PreparedStatement.class, String.class,String.class,int.class, String.class);
            insertBefore(ctClass,methodName,desc,putScr);
        } catch (NotFoundException e) {
            e.printStackTrace();
        } catch (CannotCompileException e) {
            e.printStackTrace();
        }
    }

    public static void putPreparedStatementSqlParam(PreparedStatement preparedStatement,String stmt,String method,int i,String value){
        if(!SQLPreparedStatementHook.isChangeDbSQL(stmt)){
            return;
        }
        WeakHashMap<PreparedStatement,Map<Integer, SQLPreparedParam>> preparedStatementMap = sqlParamsData.get();
        if(preparedStatementMap==null){
            preparedStatementMap = new WeakHashMap<PreparedStatement,Map<Integer, SQLPreparedParam>>();
            sqlParamsData.set(preparedStatementMap);
        }
        Map<Integer, SQLPreparedParam> integerSQLPreparedParamMap = preparedStatementMap.get(preparedStatement);
        if(integerSQLPreparedParamMap==null){
            integerSQLPreparedParamMap = new HashMap<Integer, SQLPreparedParam>();
            preparedStatementMap.put(preparedStatement,integerSQLPreparedParamMap);
        }

        SQLPreparedParam sqlPreparedParam = new SQLPreparedParam();
        sqlPreparedParam.setIndex(i);
        sqlPreparedParam.setMethod(method);
        sqlPreparedParam.setValue(value);
        integerSQLPreparedParamMap.put(i,sqlPreparedParam);
        //System.out.println("preparedStatementMap-size:"+preparedStatementMap.size());
    }


    public static Map<Integer, SQLPreparedParam> getAndRemovePreparedStatementSqlParams(PreparedStatement preparedStatement){
        Map<PreparedStatement,Map<Integer, SQLPreparedParam>> preparedStatementMap = SQLPreparedStatementHook.sqlParamsData.get();
        if(preparedStatementMap==null){
            return null;
        }
        return preparedStatementMap.remove(preparedStatement);
    }



}
