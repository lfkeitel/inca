/* global $:false, setInterval, clearInterval */

"use strict"; // jshint ignore:line
var runningRefreshTimer;
var bgRefreshTimer;

function run() {
    server.performRun(function(data) {
        updateView({"running": true, "totalDevices": 1, "finished": 0});
        progressBar.setProgressBarToRunning(true);
    });
    if (typeof runningRefreshTimer === 'undefined') {
        runningRefreshTimer = setInterval(function() { server.checkStatus(checkRequest); }, 1000);
    }
}

function checkRequest(data) {
    updateView(data);

    if (!data.running) {
        clearInterval(runningRefreshTimer);
        progressBar.setProgressBarToRunning(false);
    }
}

function updateView(data) {
    progressBar.setProgressBarActive(data.running);

    if (data.running === false) {
        $('#currentStatus').html('Idle').addClass('idleStatus');
        $('#currentStatus').removeClass('runningStatus');
    } else {
        $('#currentStatus').html('Running').addClass('runningStatus');
        $('#currentStatus').removeClass('idleStatus');
        if (typeof runningRefreshTimer === 'undefined') {
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
}

var server = {
    checkStatus: function(callback) {
        $.get('/api/status', {}, null, 'json')
            .done(function(data) {
                callback(data);
            });
    },

    performRun: function(callback) {
        $.get('/api/runnow', {}, null, 'json')
            .done(function(data) {
                callback(data);
            });
    },

    getDeviceList: function(callback) {
        $.get('/api/devicelist', {}, null, 'json')
            .done(function(data) {
                callback(data);
            });
    },
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
    server.checkStatus(updateView);
    bgRefreshTimer = setInterval(function() { server.checkStatus(updateView); }, 30000);
})();
