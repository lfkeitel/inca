/* global $:false, event */

"use strict"; // jshint ignore:line

var deviceList;

var server = {
    getDeviceList: function(callback) {
        $.get('/api/devicelist', {}, null, 'json')
            .done(function(data) {
                if (typeof callback !== 'undefined') {
                    callback(data);
                }
                return;
            });
        return;
    },
};

function check(e) {
    if (e.keyCode == 13) {
        searchList();
        e.preventDefault();
    }
}

function searchList() {
    var address = $('#searchAddress').val();
    $('#searchResults').empty();
    for (var i = 0; i < deviceList.length; i++) {
        if (deviceList[i].Address === address) {
            showSearchResult(deviceList[i]);
            break;
        }
    }
    return;
}

function showSearchResult(result) {
    $('#searchResults').append('<h4>Search Results:</h4>');


    // jshint multistr:true
    var tableHead = '<thead><tr>\
                        <td>IP Address</td>\
                        <td>Name</td>\
                        <td>Protocol</td>\
                        <td>Path</td>\
                    </tr></thead>';

    var table = '<table class="archiveList">'+tableHead+'<tr>\
                    <td>'+result.Address+'</td>\
                    <td>'+result.Name+'</td>\
                    <td>'+result.Proto+'</td>\
                    <td><a href="view/'+result.Path+'">'+result.Path+'</a></td>\
                </tr></table>';
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
    return;
})();
