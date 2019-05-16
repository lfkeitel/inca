import * as $ from "jquery";
import * as api from "./api";

function makeCallback(msg: string): (data: api.ISaveConfigResult) => void {
    return function (data: api.ISaveConfigResult) {
        if (!data.success) {
            alert(data.error);
        } else {
            alert(msg);
        }
    }
}

(function () {
    $('#saveDeviceListBtn').click(
        () => {
            const listText = $('#deviceListConfig').val() as string;
            api.saveDeviceList(listText, makeCallback("Device list saved"))
        }
    );
    $('#saveDeviceTypeBtn').click(
        () => {
            const listText = $('#deviceTypeConfig').val() as string;
            api.saveDeviceTypes(listText, makeCallback("Device type definitions saved"))
        }
    );
})();
