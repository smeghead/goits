$(function(){
  $("#columns table.table").tableDnD({
    onDragClass: 'dragging',
    dragHandle: ".cursor-grab"
  });

  $('#columns').on('click', '.edit', function(){
    var item = $(this).parents('tr').data('item');
    var win = $('#item-edit-block-tmpl').tmpl(item).addClass('modal').appendTo($('body'));
    win.on('hidden', function(){
      win.remove();
    });
    win.modal();
  });
});
// vim: set ts=2 sw=2 sts=2 expandtab fenc=utf-8:
