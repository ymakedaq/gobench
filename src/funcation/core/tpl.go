package core

const (
	tpl = `	
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<html lang="en">
<head>
  <script type="text/javascript" src="http://cdn.hcharts.cn/jquery/jquery-1.8.3.min.js"></script>
  <script type="text/javascript" src="http://cdn.hcharts.cn/highcharts/highcharts.js"></script>
  <script type="text/javascript" src="http://cdn.hcharts.cn/highcharts/exporting.js"></script>
  <script>
    $(function () {
    $('#insert').highcharts({
        chart: {
            type: 'spline'
        },
        title: {
            text: {{.XTitle.Cpu}}
        },
        xAxis: {
            categories: {{.Timeline}},
        },
        yAxis: {
            title: {
                text: {{.YTitle.Cpu}}
            },
            labels: {
                formatter: function() {
                    return this.value 
                }
            }
        }, 
        credits: {   
            text: 'ds_group@juanpi.com',
            href: 'http://www.autodb.com'  
        }, 
        tooltip: {
            crosshairs: true,
            shared: true,
            valueSuffix: ''
        },
        plotOptions: {
            spline: {
                marker: {
                    radius: 2,
                    lineColor: '#234',
                    lineWidth: 1
                }
            }
        },
        series: [{
            name: '%idle',
            data: []
	 },{
		name: 'id1',
		data: []

	 }
		]          
    });             
});     
 </script>            
</head>
<body>
  <div id="insert"></div>
</body>
</html>
	`
)
