(function() {
    function saveDevice(d, f) {
        $.post("/api/devices/save", d, f, 'json');
    }

    function prepareData() {
        var formData = $('form[name="deviceForm"]').serializeArray();
        var apiData = {};

        for (var i = 0; i < formData.length; i++) {
            var obj = formData[i];
            // Removed "device-" prefix
            var key = obj.name.substring(7);
            apiData[key] = obj.value;
        }

        if (typeof apiData.disabled === "undefined") {
            apiData.disabled = 0;
        } else if (apiData.disabled === "on") {
            apiData.disabled = 1;
        }
        return apiData;
    }

    $('#saveDeviceBtn').click(function() {
        var d = prepareData();
        saveDevice(d, function(r) {
            if (r.ErrorCode === 0) {
                alert("Device saved");
            } else {
                alert(r.ErrorMessage);
            }
        });
    });

    $('#saveListDeviceBtn').click(function() {
        var d = prepareData();
        saveDevice(d, function(r) {
            if (r.ErrorCode === 0) {
                location.assign("/devices");
            } else {
                alert(r.ErrorMessage);
            }
        });
    });
})();
