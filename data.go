/**
* @program: debug_charts
*
* @description:
*
* @author: lemo
*
* @create: 2020-01-06 21:55
**/

package debug_charts

import "html/template"

var html = template.HTMLEscapeString(`
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			<meta http-equiv="X-UA-Compatible" content="ie=edge" />
			<title>Document</title>
		</head>
		<body>
			<div><canvas id="BytesAllocatedChart"></canvas></div>
			<div><canvas id="GcPause"></canvas></div>
			<div><canvas id="Counter"></canvas></div>
		</body>
	</html>
	<script src="https://cdn.jsdelivr.net/npm/chart.js@2.8.0"></script>
	<script>
	
		let globalOption = {
			animation: { duration: 0 }, // general animation time
			hover: { animationDuration: 0 }, // duration of animations when hovering an item
			responsiveAnimationDuration: 0, // animation duration after a resize
			elements: { point: { pointStyle: "dash" } },
			scales: {
				xAxes: [],
				yAxes: [
					{
						display: true,
						ticks: {
							beginAtZero: true,
							// steps: 1000,
							stepValue: 0.1
							// max: 100
						}
						// stacked: true
						// scaleLabel: {
						// 	display: true,
						// 	labelString: "MB"
						// }
					}
				]
			},
			tooltips: {
				mode: "point"
			}
		};
	
		var BytesAllocatedChart = new Chart(document.getElementById("BytesAllocatedChart").getContext("2d"), {
			type: "line",
			data: {
				labels: [],
				datasets: [
					{
						label: "BytesAllocated",
						backgroundColor: "rgb(0, 0, 255)",
						borderColor: "rgb(0, 0, 255)",
						fill: false,
						data: []
					}
				]
			},
			options: globalOption
		});
	
		// BytesAllocatedChart.options.scales.yAxes[0].scaleLabel.labelString = "MB";
		// BytesAllocatedChart.options.scales.yAxes[0].scaleLabel.display = true;
	
		var GcPauseChart = new Chart(document.getElementById("GcPause").getContext("2d"), {
			type: "line",
			data: {
				labels: [],
				datasets: [
					{
						label: "GcPause",
						backgroundColor: "rgb(0, 99, 132)",
						borderColor: "rgb(0, 99, 132)",
						fill: false,
						data: []
					}
				]
			},
			options: globalOption
		});
	
		// GcPauseChart.options.scales.yAxes[0].scaleLabel.labelString = "MS";
		// GcPauseChart.options.scales.yAxes[0].scaleLabel.display = true;
	
		var CounterChart = new Chart(document.getElementById("Counter").getContext("2d"), {
			type: "line",
			data: {
				labels: [],
				datasets: [
					{
						label: "Goroutine",
						backgroundColor: "rgb(200, 0, 132)",
						borderColor: "rgb(200, 0, 132)",
						fill: false,
						data: []
					},
					{
						label: "ThreadCreate",
						backgroundColor: "rgb(88, 255, 99)",
						borderColor: "rgb(88, 255, 99)",
						fill: false,
						data: []
					},
					{
						label: "Heap",
						backgroundColor: "rgb(0, 0, 0)",
						borderColor: "rgb(0, 0, 0)",
						fill: false,
						data: []
					},
					{
						label: "Mutex",
						backgroundColor: "rgb(167, 98, 77)",
						borderColor: "rgb(167, 98, 77)",
						fill: false,
						data: []
					},
					{
						label: "Block",
						backgroundColor: "rgb(255, 0, 132)",
						borderColor: "rgb(255, 0, 132)",
						fill: false,
						data: []
					}
				]
			},
			options: globalOption
		});
	
		function update(info) {
			BytesAllocatedChart.data.labels.push(new Date().getSeconds());
			if (BytesAllocatedChart.data.labels.length > 120) {
				BytesAllocatedChart.data.labels.shift();
			}
	
			BytesAllocatedChart.data.datasets.forEach(dataset => {
				if (dataset.label.startsWith("BytesAllocated")) {
					var bytes = (info.BytesAllocated / 1024 / 1024).toFixed(8);
					dataset.label = "BytesAllocated: " + bytes + " MB";
					dataset.data.push(bytes);
					if (dataset.data.length > 120) {
						dataset.data.shift();
					}
				}
			});
			BytesAllocatedChart.update();
	
			GcPauseChart.data.labels.push(new Date().getSeconds());
			if (GcPauseChart.data.labels.length > 120) {
				GcPauseChart.data.labels.shift();
			}
			GcPauseChart.data.datasets.forEach(dataset => {
				if (dataset.label.startsWith("GcPause")) {
					var ms = (info.GcPause / 100000).toFixed(8);
					dataset.label = "GcPause: " + ms + " MS";
					dataset.data.push(ms);
					if (dataset.data.length > 120) {
						dataset.data.shift();
					}
				}
			});
			GcPauseChart.update();
	
			CounterChart.data.labels.push(new Date().getSeconds());
			if (CounterChart.data.labels.length > 120) {
				CounterChart.data.labels.shift();
			}
			CounterChart.data.datasets.forEach(dataset => {
				// Block
				if (dataset.label.startsWith("Block")) {
					dataset.label = "Block: " + info.Block;
					dataset.data.push(info.Block);
					if (dataset.data.length > 120) {
						dataset.data.shift();
					}
				}
	
				// Goroutine
				if (dataset.label.startsWith("Goroutine")) {
					dataset.label = "Goroutine: " + info.Goroutine;
					dataset.data.push(info.Goroutine);
					if (dataset.data.length > 120) {
						dataset.data.shift();
					}
				}
	
				// Heap
				if (dataset.label.startsWith("Heap")) {
					dataset.label = "Heap: " + info.Heap;
					dataset.data.push(info.Heap);
					if (dataset.data.length > 120) {
						dataset.data.shift();
					}
				}
	
				// Mutex
				if (dataset.label.startsWith("Mutex")) {
					dataset.label = "Mutex: " + info.Mutex;
					dataset.data.push(info.Mutex);
					if (dataset.data.length > 120) {
						dataset.data.shift();
					}
				}
	
				// ThreadCreate
				if (dataset.label.startsWith("ThreadCreate")) {
					dataset.label = "ThreadCreate: " + info.ThreadCreate;
					dataset.data.push(info.ThreadCreate);
					if (dataset.data.length > 120) {
						dataset.data.shift();
					}
				}
			});
			CounterChart.update();
		}
	
		let webSocket = new WebSocket("ws://127.0.0.1:23456");
	
		webSocket.onopen = () => {
			setInterval(() => {
				webSocket.send("");
			}, 3000);
		};
	
		webSocket.onmessage = msg => {
			let message = JSON.parse(msg.data);
			update(message.data.msg);
		};
	
		webSocket.onclose = () => {};
	
		webSocket.onerror = () => {};
	
	
	</script>

`)
