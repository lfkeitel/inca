var selectedDevices = [];

var buttonElements = {
    edit: $("#deviceEditBtn"),
    create: $("#deviceCreateBtn"),
    delete: $("#deviceDeleteBtn"),
    update: $("#deviceUpdateBtn"),
    refresh: $("#refreshBtn")
};

function goToConfigs(id) {
    event.stopPropagation();
    var url = "/configs/"+id;
    if (event.button == 1) {
        window.open(url);
    } else {
        location.assign(url);
    }
}

function checkDisabledButtons() {
    var len = selectedDevices.length;

    if (len === 0) {
        lib.disableInput(buttonElements.edit);
        lib.disableInput(buttonElements.delete);
        lib.disableInput(buttonElements.update);
    } else if (len == 1) {
        lib.enableInput(buttonElements.edit);
        lib.enableInput(buttonElements.delete);
        lib.enableInput(buttonElements.update);
    } else {
        lib.disableInput(buttonElements.edit);
        lib.enableInput(buttonElements.delete);
        lib.enableInput(buttonElements.update);
    }
}

(function() {
    $('.selectable_row').click(function() {
        var id = $(this).data("did");
        var index = $.inArray(id, selectedDevices);

        if (index > -1) {
            $('#device-'+id).prop("checked", false);
            lib.arrayRemove(selectedDevices, index);
        } else {
            $('#device-'+id).prop("checked", true);
            selectedDevices.push(id);
        }

        checkDisabledButtons();
    });
})();
