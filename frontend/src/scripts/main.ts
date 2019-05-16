import * as $ from "jquery";
import * as api from "./api";

let runningRefreshTimer: NodeJS.Timeout = null;
let bgRefreshTimer: NodeJS.Timeout = null;
const initial_view_data: api.IStatusResult = {
    stage: api.Stage.Default,
    running: true,
    totalDevices: 1,
    finished: 0
};

function run() {
    api.performRun(function (data) {
        view.updateView(initial_view_data);
        progressBar.setProgressBarToRunning(true);
    });
    if (runningRefreshTimer === null) {
        runningRefreshTimer = setInterval(() => api.checkStatus(checkRequest), 1000);
    }
}

function checkRequest(data: api.IStatusResult) {
    view.updateView(data);

    if (!data.running) {
        clearInterval(runningRefreshTimer);
        runningRefreshTimer = null;
        progressBar.setProgressBarToRunning(false);
    }
}

function stage_to_user_string(stage: string): string {
    switch (stage) {
        case "grabbing":
            return "Grabbing";
        case "loading-configuration":
            return "Loading Configuration";
        case "post-script":
            return "Post Script";
        case "pre-script":
            return "Pre Script";
    }
}

const view = {
    updateView: function (data: api.IStatusResult) {
        progressBar.setProgressBarActive(data.running);

        if (data.running === false) {
            $('#currentStatus').html('Idle').addClass('idleStatus');
            $('#currentStatus').removeClass('runningStatus');
            api.getErrorLog(view.updateLogView);
            if (bgRefreshTimer === null) {
                bgRefreshTimer = setInterval(() => api.checkStatus(view.updateView), 30000);
            }
        } else {
            if (data.stage !== '') {
                $('#currentStatus').html(`Running - ${stage_to_user_string(data.stage)}`).addClass('runningStatus');
            } else {
                $('#currentStatus').html('Running').addClass('runningStatus');
            }
            $('#currentStatus').removeClass('idleStatus');
            api.getErrorLog(view.updateLogView);
            if (runningRefreshTimer === null) {
                clearInterval(bgRefreshTimer);
                bgRefreshTimer = null;
                progressBar.setProgressBarToRunning(true);
                runningRefreshTimer = setInterval(() => api.checkStatus(checkRequest), 1000);
            }
        }

        if (data.finished < 0) {
            progressBar.setProgressBarMax(100);
            progressBar.setProgressBarValue(100);
        } else {
            progressBar.setProgressBarMax(data.totalDevices);
            progressBar.setProgressBarValue(data.finished);
        }
    },

    updateLogView: function (data: api.ErrorLine[]) {
        const table = $('<table/>');
        const tableHead = `<thead><tr>
            <td>Type</td>
            <td>Time</td>
            <td>Message</td>
            </tr></thead>`;

        table.append(tableHead);

        for (let key in data) {
            if (!data.hasOwnProperty(key))
                continue;

            const log = data[key];

            const html = `
                <tr class="${log.etype.toLowerCase()}">
                    <td class="${log.etype.toLowerCase()}">${log.etype}</td>
                    <td>${log.time}</td>
                    <td>${log.message}</td>
                </tr>
            `;

            table.append(html);
        }

        $('#appLogs').empty();
        $('#appLogs').append(table);
    }
};

function manualSingleDeviceRun() {
    const manName = $('#manName').val() as string;
    const manAddr = $('#manAddr').val() as string;
    const manType = $('#manType').val() as string;
    const manProto = $('#manProto').val() as string;
    api.runSingleDeviceGrab(
        manAddr, manType, manProto, manName,
        (data) => alert("Downloading new config. Check Status page.")
    );
}

const progressBar = {
    setProgressBarValue: function (val: number) {
        const max = parseInt($('#statusProgressBar').attr('aria-valuemax'));
        $('#statusProgressBar').css('width', ((val / max) * 100) + '%').attr('aria-valuenow', val);
    },

    setProgressBarMax: function (val: number) {
        $('#statusProgressBar').attr('aria-valuemax', val);
    },

    setProgressBarActive: function (val: boolean) {
        if (val) {
            $('#statusProgressBar').addClass('active');
            $('#statusProgressBar').addClass('progress-bar-striped');
        } else {
            $('#statusProgressBar').removeClass('active');
            $('#statusProgressBar').removeClass('progress-bar-striped');
        }
    },

    setProgressBarToRunning: function (val: boolean) {
        if (val) {
            $('#statusProgressBar').addClass('progress-bar-danger');
            $('#statusProgressBar').removeClass('progress-bar-success');
        } else {
            $('#statusProgressBar').removeClass('progress-bar-danger');
            $('#statusProgressBar').addClass('progress-bar-success');
        }
    },
}; // progressBar

(function () {
    $('#startArchiveBtn').click(run);
    $('#manualRunBtn').click(manualSingleDeviceRun);
    api.checkStatus(view.updateView);
    api.getErrorLog(view.updateLogView);
    bgRefreshTimer = setInterval(function () { api.checkStatus(view.updateView); }, 30000);
})();
