(function() {
    var buttons = {
        cpView: $('#cpViewBtn'),
        dtView: $('#dtViewBtn'),
    }
    var partialDiv = $('#partial-view');

    var selectedProfile = "";

    buttons.cpView.click(function(){
        buttons.dtView.removeClass("active");
        buttons.cpView.addClass("active");
        getPartial("cp", cpPrepare);
    });

    buttons.dtView.click(function(){
        buttons.cpView.removeClass("active");
        buttons.dtView.addClass("active");
        getPartial("dt", dtPrepare);
    });

    function getPartial(p, prepare)  {
        $.get("/admin/partial/"+p, {}, function(r) {
            partialDiv.html(r);
            prepare();
        }, "html");
    }

    function cpPrepare() {
        $('#cp-ids').change(function() {
            var id = $(this).val();
            $('.cp-form').hide();
            $('#cp-'+id+'-form').show();
            selectedProfile = id;
        });

        $('.save-btn').click(function() {
            var d = prepateCPForm();
            saveCP(d, function(r) {
                if (r.ErrorCode === 0) {
                    getPartial("cp", cpPrepare);
                } else {
                    alert(r.ErrorMessage);
                }
            });
        });

        $('.delete-btn').click(function() {
            $.post("/api/cp/delete", {ids: selectedProfile}, function(r) {
                if (r.ErrorCode === 0) {
                    getPartial("cp", cpPrepare);
                } else {
                    alert(r.ErrorMessage);
                }
            }, 'json');
        });
    }

    function prepateCPForm() {
        var formData = $('form[id="cp-'+selectedProfile+'-form"]').serializeArray();
        var apiData = {};
        apiData.id = selectedProfile;

        for (var i = 0; i < formData.length; i++) {
            var obj = formData[i];
            // Removed "cp-" prefix
            var key = obj.name.substring(3);
            apiData[key] = obj.value;
        }

        return apiData;
    }

    function saveCP(d, f) {
        $.post("/api/cp/save", d, f, 'json');
    }

    function dtPrepare() {
        console.log("loaded2");
    }
    getPartial("cp", cpPrepare);
})();
