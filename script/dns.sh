#! /bin/bash
#自定义你们公司获取域名对应ip列表的脚本
#输出形式如下,idc和ip直接用逗号隔开
#idc1,220.181.38.148
#idc1,39.156.69.79

domain=$1
for i in `nslookup $domain|grep Address|grep -v 53|awk '{print $2}'`;
	do echo "idc1",$i;
done;
