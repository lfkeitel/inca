import * as api from "./api";
import * as $ from "jquery";

let deviceList: api.Device[];

function singleRun(address: string, brand: string, proto: string, name: string) {
    api.runSingleDeviceGrab(address, brand, proto, name,
        () => alert("Downloading new config. Check Status page."));
}

function check(e: JQuery.KeyPressEvent) {
    if (e.keyCode == 13) {
        searchList();
        e.preventDefault();
    }
}

function searchList() {
    const query = ($('#searchAddress').val() as string).toLowerCase();
    $('#searchResults').empty();
    const results = deviceList.filter((device) =>
        device.address.includes(query) || device.name.toLowerCase().includes(query));
    showSearchResult(results);
}

function showSearchResult(results: api.Device[]) {
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
            <td>${result.address}</td>
            <td>${result.name}</td>
            <td>${result.proto}</td>
            <td><a href="view/${result.path}">${result.path}</a></td>
        </tr>`);
    });

    $('#searchResults').append(table);
    return;
}

(function () {
    $('#searchBtn').click(searchList);
    $('#searchAddress').keypress((event) => check(event));
    api.getDeviceList((data) => deviceList = data.devices);

    $('.single-run-btn').click(function () {
        const address = this.dataset.address;
        const brand = this.dataset.man;
        const name = this.dataset.name;
        const proto = this.dataset.proto;
        singleRun(address, brand, proto, name);
    });
})();
