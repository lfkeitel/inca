(function() {
    var buttons = {
        cpView: $('#cpViewBtn'),
        dtView: $('#dtViewBtn'),
    }
    var partialDiv = $('#partial-view');

    buttons.cpView.click(function(){
        buttons.dtView.removeClass("active");
        buttons.cpView.addClass("active");
        getPartial("cp");
    });

    buttons.dtView.click(function(){
        buttons.cpView.removeClass("active");
        buttons.dtView.addClass("active");
        getPartial("dt");
    });

    function getPartial(p)  {
        $.get("/admin/partial/"+p, {}, function(r) {
            partialDiv.html(r);
        }, "html");
    }

    getPartial("cp");
})();
