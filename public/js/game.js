function initMainPanel($) {
    var mainPanels = $('.mainPanel');
    function switchMainPanel(name) {
        mainPanels.hide();
        var m = this.id.match(/^(.+?)NavItem$/);
        $('#' + m[1]).show();
        return false;
    }
    $('.mainPanelNavItem').click(switchMainPanel);
    $('.mainPanelNavItem.default').click();
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
    $('.editPanelNavItem.default').click();
}
function finishInit($) {
    // TODO: 色々と待つ処理
    // setTimeout はダミー的処理
    setTimeout(function () {$('#loading').hide(); }, 100);
}
jQuery(initMainPanel);
jQuery(initEditPanel);
jQuery(finishInit);
