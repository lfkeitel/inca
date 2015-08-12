(function() {
    // Keeps track of the device ids for selected rows in table
    var selectedDevices = [];

    // jQuery objects of buttons
    var buttonElements = {
        edit: $("#deviceEditBtn"),
        create: $("#deviceCreateBtn"),
        delete: $("#deviceDeleteBtn"),
        update: $("#deviceUpdateBtn"),
        refresh: $("#refreshBtn"),
        search: $("#searchBtn"),
    };

    // Collection of click handlers
    var clickHandlers = {
        buttons: {
            edit: function() {
                if (!selectedDevices[0]) {
                    return;
                }

                var url = "/devices/"+selectedDevices[0];
                location.assign(url);
            },

            create: function() {
                location.assign("/devices/new");
            },

            delete: function() {
                $.post('/api/devices/delete', {deviceids: JSON.stringify(selectedDevices)}, function(data) {
                    if (data.ErrorCode === 0) {
                        location.reload();
                    } else {
                        alert(data.ErrorMessage);
                    }
                }, "json");
            },

            update: function() {
                $.post('/api/devices/update', {deviceids: JSON.stringify(selectedDevices)}, function(data) {
                    if (data.ErrorCode === 0) {
                        alert("Device status and configurations are updating");
                    } else {
                        alert(data.ErrorMessage);
                    }
                }, "json");
            },

            refresh: function() {
                location.reload();
            },

            search: function() {
                var query = $('#searchbox').val();
                location.assign("/devices/?query="+encodeURIComponent(query));
            },
        },

        configs: function(id, button) {
            var url = "/configs/"+id;
            if (button === 1) {
                window.open(url);
            } else {
                location.assign(url);
            }
        },
    };

    // Display the appropiate buttons given the number of selected rows
    function checkDisabledButtons() {
        var len = selectedDevices.length;

        if (len === 0) {
            $.disableInput(buttonElements.edit);
            $.disableInput(buttonElements.delete);
            $.disableInput(buttonElements.update);
        } else if (len == 1) {
            $.enableInput(buttonElements.edit);
            $.enableInput(buttonElements.delete);
            $.enableInput(buttonElements.update);
        } else {
            $.disableInput(buttonElements.edit);
            $.enableInput(buttonElements.delete);
            $.enableInput(buttonElements.update);
        }
    }

    // Wire up config icons
    $('.config-clickable').click(function(e) {
        e.stopPropagation();
        var me = $(this);
        var id = me.data("configid");
        clickHandlers.configs(id, e.button);
    });

    // Wire up row selection
    $('.selectable_row').click(function() {
        var id = $(this).data("did");
        var index = $.inArray(id, selectedDevices);

        if (index > -1) {
            $('#device-'+id).prop("checked", false);
            $.arrayRemove(selectedDevices, index);
        } else {
            $('#device-'+id).prop("checked", true);
            selectedDevices.push(id);
        }

        checkDisabledButtons();
    });

    // Wire up the buttons
    $.each(buttonElements, function(key, btnEl) {
        var handler = clickHandlers.buttons[key];
        if (handler) {
            btnEl.click(handler);
        }
    });

    $('#searchbox').keypress(function(e) {
        if (e.keyCode === 13) {
            e.stopPropagation();
            clickHandlers.buttons.search();
            return false;
        }
    })
})();
