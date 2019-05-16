/* global $:false, alert */

"use strict"; // jshint ignore:line

function saveDeviceList() {
    const listText = $('#deviceListConfig').val();

    $.post('/api/savedevicelist', { text: encodeURIComponent(listText) }, null, "json")
        .done(function(data) {
            if (!data.success) {
                alert(data.error);
            } else {
                alert("Device list saved");
            }
        });
}

function saveDeviceTypes() {
    const listText = $('#deviceTypeConfig').val();

    $.post('/api/savedevicetypes', { text: encodeURIComponent(listText) }, null, "json")
        .done(function(data) {
            if (!data.success) {
                alert(data.error);
            } else {
                alert("Device type definitions saved");
            }
        });

}
