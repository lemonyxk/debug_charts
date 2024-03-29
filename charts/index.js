let globalOption = {
    animation: {duration: 0}, // general animation time
    hover: {animationDuration: 0}, // duration of animations when hovering an item
    responsiveAnimationDuration: 0, // animation duration after a resize
    elements: {point: {pointStyle: "dash"}},
    scales: {
        xAxes: [
            {
                display: false
            }
        ],
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
        mode: "nearest",
        // intersect: false,
        position: "nearest"
    }
};

var BytesAllocatedChart = new Chart(document.getElementById("BytesAllocatedChart").getContext("2d"), {
    type: "line",
    data: {
        labels: [],
        datasets: [
            {
                label: "BytesAllocated",
                backgroundColor: "#c901f1",
                borderColor: "#c901f1",
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
                backgroundColor: "#ed1941",
                borderColor: "#ed1941",
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
                backgroundColor: "#53261f",
                borderColor: "#53261f",
                fill: false,
                data: []
            },
            {
                label: "ThreadCreate",
                backgroundColor: "#fdb933",
                borderColor: "#fdb933",
                fill: false,
                data: []
            },
            {
                label: "Heap",
                backgroundColor: "#00a6ac",
                borderColor: "#00a6ac",
                fill: false,
                data: []
            },
            {
                label: "Mutex",
                backgroundColor: "#999d9c",
                borderColor: "#999d9c",
                fill: false,
                data: []
            },
            {
                label: "Block",
                backgroundColor: "#fedcbd",
                borderColor: "#fedcbd",
                fill: false,
                data: []
            }
        ]
    },
    options: globalOption
});

const maxColumn = 3600;

var sec = 0;

function update(info) {
    var second = new Date().getSeconds();

    BytesAllocatedChart.data.labels.push(sec === second ? sec : "");
    if (BytesAllocatedChart.data.labels.length > maxColumn) {
        BytesAllocatedChart.data.labels.shift();
    }

    BytesAllocatedChart.data.datasets.forEach(dataset => {
        if (dataset.label.startsWith("BytesAllocated")) {
            var bytes = (info.BytesAllocated / 1024 / 1024).toFixed(8);
            dataset.label = "BytesAllocated: " + bytes + " MB";
            dataset.data.push(bytes);
            if (dataset.data.length > maxColumn) {
                dataset.data.shift();
            }
        }
    });
    BytesAllocatedChart.update();

    GcPauseChart.data.labels.push(sec === second ? sec : "");
    if (GcPauseChart.data.labels.length > maxColumn) {
        GcPauseChart.data.labels.shift();
    }
    GcPauseChart.data.datasets.forEach(dataset => {
        if (dataset.label.startsWith("GcPause")) {
            var ms = (info.GcPause / 1e6).toFixed(8);
            dataset.label = "GcPause: " + ms + " MS";
            dataset.data.push(ms);
            if (dataset.data.length > maxColumn) {
                dataset.data.shift();
            }
        }
    });
    GcPauseChart.update();

    CounterChart.data.labels.push(sec === second ? sec : "");
    if (CounterChart.data.labels.length > maxColumn) {
        CounterChart.data.labels.shift();
    }
    CounterChart.data.datasets.forEach(dataset => {
        // Block
        if (dataset.label.startsWith("Block")) {
            dataset.label = "Block: " + info.Block;
            dataset.data.push(info.Block);
            if (dataset.data.length > maxColumn) {
                dataset.data.shift();
            }
        }

        // Goroutine
        if (dataset.label.startsWith("Goroutine")) {
            dataset.label = "Goroutine: " + info.Goroutine;
            dataset.data.push(info.Goroutine);
            if (dataset.data.length > maxColumn) {
                dataset.data.shift();
            }
        }

        // Heap
        if (dataset.label.startsWith("Heap")) {
            dataset.label = "Heap: " + info.Heap;
            dataset.data.push(info.Heap);
            if (dataset.data.length > maxColumn) {
                dataset.data.shift();
            }
        }

        // Mutex
        if (dataset.label.startsWith("Mutex")) {
            dataset.label = "Mutex: " + info.Mutex;
            dataset.data.push(info.Mutex);
            if (dataset.data.length > maxColumn) {
                dataset.data.shift();
            }
        }

        // ThreadCreate
        if (dataset.label.startsWith("ThreadCreate")) {
            dataset.label = "ThreadCreate: " + info.ThreadCreate;
            dataset.data.push(info.ThreadCreate);
            if (dataset.data.length > maxColumn) {
                dataset.data.shift();
            }
        }
    });
    CounterChart.update();

    if (sec != second) {
        sec = second;
    }
}

function ws() {

    let protocol = 'ws';

    if (window.location.protocol === 'https:') protocol = 'wss';

    let webSocket = new WebSocket(protocol + "://" + window.location.hostname + ":" + window.location.port + "/debug/feed/");

    webSocket.onopen = () => {
        webSocket.send(JSON.stringify({"event": "/debug/login"}))
        setInterval(() => {
            webSocket && webSocket.send("");
        }, 3000);
    };

	webSocket.onmessage = async (msg) => {
		var text = await (new Response(msg.data)).text();
        let message = JSON.parse(text.slice(22));
        message.msg = message.msg || [];
        for (let i = 0; i < message.msg.length; i++) {
            update(message.msg[i]);
        }
    };

    webSocket.onclose = () => {
        webSocket = null;
        setTimeout(() => {
            ws();
        }, 1000);
    };

    webSocket.onerror = () => {
    };
}

ws();