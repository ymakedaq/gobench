package datahandle

const (
	basetpl = `
	<html>
		<head>
		<meta charset="utf-8"><link rel="icon" href="https://static.jianshukeji.com/highcharts/images/favicon.ico">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style></style>
		<script src="http://cdn.hcharts.cn/jquery/jquery-1.8.3.min.js"></script>
		<script src="https://img.hcharts.cn/highcharts/highcharts.js"></script>
		<script	src="https://img.hcharts.cn/highcharts/modules/exporting.js"></script>
		<script src="https://img.hcharts.cn/highcharts-plugins/highcharts-zh_CN.js"><script>
		</head>
	<body>
		<div id="container" style="min-width:400px;height:400px"></div>
		<script>
		{{.Js_code}}
		</script>
	</body>
	<html>
	`
)
