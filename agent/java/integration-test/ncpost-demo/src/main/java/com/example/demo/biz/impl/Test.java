package com.example.demo.biz.impl;

import java.util.StringTokenizer;

public class Test {
	public static void main(String[] args) {
		String str = "100|66,55:200|567,90:102|43,54";

		StringTokenizer strToke = new StringTokenizer(str, ":,|");// 默认不打印分隔符
		// StringTokenizer strToke=new StringTokenizer(str,":,|",true);//打印分隔符
		// StringTokenizer strToke=new StringTokenizer(str,":,|",false);//不打印分隔符
		while(strToke.hasMoreTokens()){
		    System.out.println(strToke.nextToken());
		}
	}
}
