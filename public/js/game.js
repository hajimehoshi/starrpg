function initMainPanel($) {
    var mainPanels = $('.mainPanel');
    function switchMainPanel(name) {
        mainPanels.hide();
        var m = this.id.match(/^(.+?)NavItem$/);
        $('#' + m[1]).show();
        return false;
    }
    $('.mainPanelNavItem').click(switchMainPanel);
}
function initEditPanel($) {
    var editPanels = $('.editPanel');
    function switchEditPanel(name) {
        editPanels.hide();
        var m = this.id.match(/^(.+?)NavItem$/);
        $('#' + m[1]).show();
        return false;
    }
    $('.editPanelNavItem').click(switchEditPanel);
}
jQuery(initMainPanel);
jQuery(initEditPanel);
