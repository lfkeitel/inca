/* global $:false, event, alert */

"use strict"; // jshint ignore:line

let deviceList;

if (typeof server === 'undefined') {
    var server = {};
}

server.getDeviceList = function(callback) {
    $.get('/api/devicelist', {}, null, 'json')
        .done(function(data) {
            if (typeof callback !== 'undefined') {
                callback(data);
            }
            return;
        });
    return;
};

server.runSingleDeviceGrab = function(address, brand, proto, name, callback) {
    $.get('/api/singlerun', { hostname: address, name: name, proto: proto, brand: brand }, null, 'json')
        .done(function(data) {
            if (typeof callback !== 'undefined') {
                callback(data);
            }
            return;
        });
    return;
};

function singleRun(address, brand, proto, name) {
    server.runSingleDeviceGrab(address, brand, proto, name, function(data) { alert("Downloading new config. Check Status page."); });
    return;
}

function check(e) {
    if (e.keyCode == 13) {
        searchList();
        e.preventDefault();
    }
}

function searchList() {
    const query = $('#searchAddress').val().toLowerCase();
    $('#searchResults').empty();
    const results = deviceList.filter((device) =>
        device.Address.includes(query) || device.Name.toLowerCase().includes(query));
    showSearchResult(results);
}

function showSearchResult(results) {
    $('#searchResults').append('<h4>Search Results:</h4>');

    const table = $(`
    <table class="archiveList">
        <thead>
            <tr>
                <td>IP Address</td>
                <td>Name</td>
                <td>Protocol</td>
                <td>Path</td>
            </tr>
        </thead>
    </table>`);

    results.forEach(result => {
        table.append(`<tr>
            <td>${result.Address}</td>
            <td>${result.Name}</td>
            <td>${result.Proto}</td>
            <td><a href="view/${result.Path}">${result.Path}</a></td>
        </tr>`);
    });

    $('#searchResults').append(table);
    return;
}

(function() {
    $('#searchBtn').click(searchList);
    $('#searchAddress').keypress(function() { check(event); });
    server.getDeviceList(function(data) {
        deviceList = data.Devices;
        return;
    });
    $('.single-run-btn').click(function() {
        const address = this.dataset.address;
        const brand = this.dataset.man;
        const name = this.dataset.name;
        const proto = this.dataset.proto;
        singleRun(address, brand, proto, name);
    });
})();
