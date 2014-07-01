package gottp

import (
    "time"
)

const layout = "Jan 2, 2006 at 3:04pm (MST)"

func ErrorTemplate(stack string, traceback string) string {
    return "<html>"+
    "<body style=\"background-color: #f6eef6;\">"+
    "<div style=\"padding: 10px; font-family: Lucida Grande, Lucida, Verdana, sans-serif; font-size:12px; max-width: 540px; margin: 0 auto;\">"+
    "<div style=\"padding-top: 0px; padding-bottom: 10px; padding-right: 10px; padding-left: 10px;\">"+
    "<div style='padding: 0; text-align: right;'>"+
    "<h1 style='color:#440000; padding: 0; margin-top: 10px; margin-bottom: 10px;font-size: 12px;'>"+
    time.Now().Format(layout)+
    "</h1>"+
    "</div>"+
    "<div style=\"padding: 10px; margin-top: 10px; margin-bottom: 10px; font-size: 12px; background-color:#FFFDCF;\">"+
    stack+
    "</div>"+
    "<div style='background-color: #eeeeee; margin-top: 10px; margin-bottom: 10px;border: 1px #440000 solid; padding-bottom: 10px;'>"+
    "<h1 style='font-size: 12px; padding: 10px; margin: 0; background-color: #440000; color: #ffffff; text-align: right'>Traceback</h1>"+
    "<div style='padding-left: 10px; padding-right: 10px'>"+
    "<pre style='background-color:#F2F2F2; padding: 10px; margin:4px;'> "+traceback+" </pre>"+
    "</div>"+
    "</div>"+
    "<div style=\"padding: 2px; font-size:11px;\">"+
    "<span><a style=\"text-decoration: none; font-weight: bold;\" href=\"http://www.siminars.com\">siminars.com</a></span>"+
    "<span style=\"float:right\"><a style=\"text-decoration: none; font-size: 9px\" href=\"http://siminars.com/terms.sv#copyright\">&copy; 2009-2012</span>"+
    "</div>"+
    "</div>"+
    "</div>"+
    "</body>"+
    "</html>"
}
