$.extend({
    // Array Remove - By John Resig (MIT Licensed)
    arrayRemove: function(array, from, to) {
        var rest = array.slice((to || from) + 1 || array.length);
        array.length = from < 0 ? array.length + from : from;
        return array.push.apply(array, rest);
    },

    disableInput: function(elem) {
        $.changeDisabledState(elem, true);
    },

    enableInput: function(elem) {
        $.changeDisabledState(elem, false);
    },

    changeDisabledState: function(elem, disabled) {
        elem = (elem instanceof jQuery) ? elem : $(elem);
        elem.prop("disabled", disabled);
    }
});
