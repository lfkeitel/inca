(function() {

var statusPieChart;

function updatePieChartDownDevices(up, down, unknown) {
    down = (typeof down === "undefined") ? 0 : down;
    unknown = (typeof unknown === "undefined") ? 0 : unknown;

    if (statusPieChart) {
        statusPieChart.destroy();
    }
    statusPieChart = new d3pie("upDownPie", {
        size: {
            canvasHeight: 250,
            canvasWidth: 400
        },
        data: {
            content: [
                {label: "Down", value: down, color: "#FF0000"},
                {label: "Up", value: up, color: "#00FF00"},
                {label: "Unknown", value: unknown, color: "#0000FF"}
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
        }
    });
}

function goToConfigs(id, button) {
    var url = "/configs/"+id;
    if (button === 1) {
        window.open(url);
    } else {
        location.assign(url);
    }
}

function updateDeviceStatus(timeout) {
    $.getJSON("/api/devices/status", {}, function(r) {
        if (r.ErrorCode === 0) {
            updatePieChartDownDevices(r.Data.Up, r.Data.Down, r.Data.Unknown);

            if (r.Data.Down) {
                var downT = $("<tbody>").attr("id", "down-table-body");

                for (var i = 0; i < r.Data.DownDevices.length; i++) {
                    var device = r.Data.DownDevices[i];
                    var html = "<tr>"+
                        "<td>"+device.Name+"</td>";
                    if (device.Status.Status === 2) {
                        html += "<td>Offline</td>";
                    } else {
                        html += "<td>Unknown</td>";
                    }

                    if (device.Status.LastPolled > 0) {
                        var date = new Date(device.Status.LastPolled*1000);
                        html += "<td>"+
                            date.toLocaleDateString()+" "+
                            date.toLocaleTimeString()+"</td>";
                    } else {
                        html += "<td>Never</td>";
                    }

                    html += "</tr>"
                    var row = $(html).click(function(e) {
                        goToConfigs(device.DeviceID, e.button);
                    });
                    downT.append(row);
                }

                $("#down-table-body").replaceWith(downT);
            } else {
                var downT = $("<tbody>").attr("id", "down-table-body");
                downT.append("<tr>" +
                    "<td>No down devices</td>" +
                    "<td>&nbsp;</td>" +
                    "<td>&nbsp;</td>" +
                "</tr>");
                $("#down-table-body").replaceWith(downT);
            }
        } else {
            $('#upDownPie').html("<h4>Error generating chart</h4>");
        }
    });

    if (timeout) {
        setTimeout(function() {
            updateDeviceStatus(timeout);
        }, timeout);
    }
}

updateDeviceStatus(300000);

})();
