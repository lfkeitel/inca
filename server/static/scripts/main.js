/* global $:false, setInterval, clearInterval, alert */

"use strict"; // jshint ignore:line
var runningRefreshTimer = null;
var bgRefreshTimer = null;

function run() {
    server.performRun(function(data) {
        view.updateView({"running": true, "totalDevices": 1, "finished": 0});
        progressBar.setProgressBarToRunning(true);
    });
    if (runningRefreshTimer === null) {
        runningRefreshTimer = setInterval(function() { server.checkStatus(checkRequest); }, 1000);
    }
}

function checkRequest(data) {
    view.updateView(data);

    if (!data.running) {
        clearInterval(runningRefreshTimer);
        runningRefreshTimer = null;
        progressBar.setProgressBarToRunning(false);
    }
}

var view = {
    updateView: function(data) {
        progressBar.setProgressBarActive(data.running);

        if (data.running === false) {
            $('#currentStatus').html('Idle').addClass('idleStatus');
            $('#currentStatus').removeClass('runningStatus');
            server.getErrorLog(view.updateLogView);
            if (bgRefreshTimer === null) {
                bgRefreshTimer = setInterval(function() { server.checkStatus(view.updateView); }, 30000);
            }
        } else {
            $('#currentStatus').html('Running').addClass('runningStatus');
            $('#currentStatus').removeClass('idleStatus');
            server.getErrorLog(view.updateLogView);
            if (runningRefreshTimer === null) {
                clearInterval(bgRefreshTimer);
                bgRefreshTimer = null;
                progressBar.setProgressBarToRunning(true);
                runningRefreshTimer = setInterval(function() { server.checkStatus(checkRequest); }, 1000);
            }
        }

        if (data.finished < 0) {
            progressBar.setProgressBarMax(100);
            progressBar.setProgressBarValue(100);
        } else {
            progressBar.setProgressBarMax(data.totalDevices);
            progressBar.setProgressBarValue(data.finished);
        }

        return;
    },

    updateLogView: function(data) {
        // jshint multistr:true
        var table = $('<table/>');
        var tableHead = '<thead><tr>\
            <td>Type</td>\
            <td>Time</td>\
            <td>Message</td>\
            </tr></thead>';

        table.append(tableHead);

        for (var key in data) {
            if (!data.hasOwnProperty(key))
                continue;

            var log = data[key];

            var html = '<tr class="'+log.Etype.toLowerCase()+'">\
                <td class="'+log.Etype.toLowerCase()+'">'+log.Etype+'</td>\
                <td>'+log.Time+'</td>\
                <td>'+log.Message+'</td>\
                </tr>';

            table.append(html);
        }

        $('#appLogs').empty();
        $('#appLogs').append(table);
        return;
    }
};

function manualSingleDeviceRun() {
    var manName = $('#manName').val();
    var manAddr = $('#manAddr').val();
    var manType = $('#manType').val();
    var manProto = $('#manProto').val();
    server.runSingleDeviceGrab(manAddr, manType, manProto, manName, function(data){alert("Downloading new config. Check Status page.");});
    return;
}

var server = {
    checkStatus: function(callback) {
        $.get('/api/status', {}, null, 'json')
            .done(function(data) {
                if (typeof callback !== 'undefined') {
                    callback(data);
                }
                return;
            });
    },

    performRun: function(callback) {
        $.get('/api/runnow', {}, null, 'json')
            .done(function(data) {
                if (typeof callback !== 'undefined') {
                    callback(data);
                }
                return;
            });
    },

    getDeviceList: function(callback) {
        $.get('/api/devicelist', {}, null, 'json')
            .done(function(data) {
                if (typeof callback !== 'undefined') {
                    callback(data);
                }
                return;
            });
    },

    getErrorLog: function(callback) {
        $.get('/api/errorlog', {limit: 10}, null, 'json')
            .done(function(data) {
                if (typeof callback !== 'undefined') {
                    callback(data);
                }
                return;
            });

    }
}; // server

var progressBar = {
    setProgressBarValue: function(val) {
        var max = $('#statusProgressBar').attr('aria-valuemax');
        $('#statusProgressBar').css('width', ((val/max)*100)+'%').attr('aria-valuenow', val);
        return;
    },

    setProgressBarMax: function(val) {
        $('#statusProgressBar').attr('aria-valuemax', val);
        return;
    },

    setProgressBarActive: function(val) {
        if (val) {
            $('#statusProgressBar').addClass('active');
            $('#statusProgressBar').addClass('progress-bar-striped');
        } else {
            $('#statusProgressBar').removeClass('active');
            $('#statusProgressBar').removeClass('progress-bar-striped');
        }
    },

    setProgressBarToRunning: function(val) {
        if (val) {
            $('#statusProgressBar').addClass('progress-bar-danger');
            $('#statusProgressBar').removeClass('progress-bar-success');
        } else {
            $('#statusProgressBar').removeClass('progress-bar-danger');
            $('#statusProgressBar').addClass('progress-bar-success');
        }
    },
}; // progressBar

(function() {
    $('#startArchiveBtn').click(run);
    $('#manualRunBtn').click(manualSingleDeviceRun);
    server.checkStatus(view.updateView);
    server.getErrorLog(view.updateLogView);
    bgRefreshTimer = setInterval(function() { server.checkStatus(view.updateView); }, 30000);
})();
