$(function(){
  $('#condition-display-toggle').click(function(){
    var button = $(this);
    button.text($(this).is(':visible') ? '詳細を隠す' : '詳細を開く');
    $('.detail-condition').toggle();
  });
});
// vim: set ts=2 sw=2 sts=2 expandtab fenc=utf-8:
