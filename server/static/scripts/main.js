function updatePieChartDownDevices(up, down) {
    var pie = new d3pie("upDownPie", {
        size: {
            canvasHeight: 250,
            canvasWidth: 350
        },
        data: {
            content: [
                {label: "Down", value: down, color: "#FF0000"},
                {label: "Up", value: up, color: "#00FF00"}
            ]
        },
        labels: {
            outer: {
                pieDistance: 20
            },
            mainLabel: {
                fontSize: 20,
                color: "#000000"
            },
            lines: {
                color: "#000000"
            },
            percentage: {
                fontSize: 0
            }
        },
        tooltips: {
            enabled: true,
            string: "Devices {label}: {value}",
            styles: {
                backgroundOpacity: 0.8,
                fontSize: 15
            }
        },
        misc: {
            pieCenterOffset: {
                x: -50
            }
        }
    });
}

updatePieChartDownDevices(47, 5);
