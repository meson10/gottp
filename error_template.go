package gottp

const templateContent = `<html>
<body style='background-color: #f6eef6;'>
<div style='padding: 10px; font-size:12px; max-width: 540px; margin: 0 auto;'>
<div style='padding-top: 0px; padding-bottom: 10px; padding-right: 10px; padding-left: 10px;'>

<div style='padding: 0; text-align: right;'>
<h1 style='color:#440000; padding: 0; margin-top: 10px; margin-bottom: 10px;font-size: 12px;'>
{{.Timestamp}}
</h1>
</div>

<div style='padding: 10px; margin-top: 10px; margin-bottom: 10px; font-size: 12px; background-color:#FFFDCF;'>
<b>Host</b>: {{.Host}} <br/>
<b>Method</b>: {{.Method}} <br/>
<b>Protocol</b>: {{.Protocol}} <br/>
<b>RemoteIP</b>: {{.RemoteIP}} <br/>
<b>URI</b>: {{.URI}} <br/>
<b>Referer</b>: {{.Referer}} <br/>
<b>Arguments</b>: {{.Arguments}} <br/>
</div>

<div style='background-color: #eeeeee; margin-top: 10px; margin-bottom: 10px;border: 1px #440000 solid; padding-bottom: 10px;'>
<h1 style='font-size: 12px; padding: 10px; margin: 0; background-color: #440000; color: #ffffff; text-align: right'>
Traceback
</h1>

<div style='padding-left: 10px; padding-right: 10px'>
<pre style='background-color:#F2F2F2; padding: 10px; margin:4px;'>
{{.Traceback}}
</pre>
</div>
</div>

</div>
</div>
</body>
</html>
`
