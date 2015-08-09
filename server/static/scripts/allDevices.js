function goToConfigs(id) {
    var url = "/configs/"+id;
    if (event.button == 1) {
        window.open(url);
    } else {
        location.assign(url);
    }
}
