// display functions
function updateTime() {
  var d = new Date();
  var text = padTwoDigit(d.getDate())+'.'+padTwoDigit(d.getMonth()+1)+'.'+d.getFullYear()+' '+padTwoDigit(d.getHours())+':'+padTwoDigit(d.getMinutes())+':'+padTwoDigit(d.getSeconds());
  $("#time").text(text);
}

function padTwoDigit(input) {
  if (input < 10) {
    return "0" + input;
  }
  return input;
}

//Entry functions
function submitAddEntry(){
  $.post("/addEntry",$("#add_entry_form").serialize())
    .done(function(){
      loadLogs();
    });
  hideEntryDialog();
}

function loadLogs() {
  $.post("/loadLogs",$("#from_to_form").serialize())
    .done(function(data){
      $("tbody").html(data)
    });

  $.post("/stats",$("#from_to_form").serialize())
    .done(function(data){
      $("#stats-block").html(data)
    });

  var f=$("#from_date").val();
  var t=$("#to_date").val();
  $("#pdf-btn").attr("href","/pdf?dateFrom="+f+"&dateTo="+t);
}

function editEntry(elem) {
  var id = elem.dataset.id
  $.get("/edit?id="+id)
    .done(function(data) {
      var entry = JSON.parse(data);
      if (entry.Type === "Work time") {
        $("#begin").val(entry.DateFrom+" "+entry.TimeFrom);
        $("#end").val(entry.DateTo+" "+entry.TimeTo);
      } else {
        $("#begin").val(entry.DateFrom);
        $("#end").val(entry.DateTo);
      }
      $("#type").val(entry.Type);
      $("#create_type").val("update");
      $("#entry_id").val(id);

      $("#add_entry_dialog").modal("toggle");
    });
}

function deleteEntry(elem) {
  var id = elem.dataset.id
  $.get("/delete?id="+id)
    .done(function() {
      loadLogs();
    });

}

function activeEntry() {
  var id = document.getElementById("active_entry").dataset.active;
  $.get("/activeEntry?id="+id)
    .done(function(data) {
      loadLogs();
      document.getElementById("active_entry").dataset.active=data;
      if (id === '') {
        $("#active_entry").text("Stop");
        $("#active_entry").removeClass("btn-success");
        $("#active_entry").addClass("btn-danger");
      } else {
        $("#active_entry").text("Start");
        $("#active_entry").removeClass("btn-danger");
        $("#active_entry").addClass("btn-success");
      }

    });
}

//user functions
function loadUsers() {
  $.get("/loadUsers")
    .done(function(data){
      $("tbody").html(data)
    });
}

function addUser() {
  cleanDialog();
  $("#type").val('create');
  $("#username").val('');
  $("#name").val('');
  $("#user_id").val('');
  $("#notRound").prop('checked',false);
  $("#edit_user_dialog").modal("toggle");
}

function editUser(elem) {
  cleanDialog();
  var id = elem.dataset.id
  $.get("/editUser?id="+id)
    .done(function(data) {
      var entry = JSON.parse(data);
      $("#username").val(entry.Username);
      $("#name").val(entry.Name);
      $("#user_id").val(id);
      $("#type").val('edit');

      $("#edit_user_dialog").modal("toggle");
    });
}

function disableUser(elem) {
  var id = elem.dataset.id;
  $.get("/toggleUser?action=disable&id="+id)
  .done(function() {
    loadUsers();
  });
}

function enableUser(elem) {
  var id = elem.dataset.id;
  $.get("/toggleUser?action=enable&id="+id)
  .done(function() {
    loadUsers();
  });
}

function submitUser(){
  var pwd = $("#password").val();
  var rep = $("#repeat").val();
  if (pwd !== rep) {
    $("#pwd_group").addClass("has-error");
    $("#rep_group").addClass("has-error");
    $("#password-alert").show();
    return;
  }
  $.post("/storeUser",$("#user_edit_form").serialize())
    .done(function(){
      loadUsers();
    });
  $("#edit_user_dialog").modal("toggle");
}

function cleanDialog() {
  $("#pwd_group").removeClass("has-error");
  $("#rep_group").removeClass("has-error");
  $("#password-alert").hide();
  $("#password").val('');
  $("#repeat").val('');
}

function updateTotal() {
  var inputs = $('.day-input');
  var total=0.0;

  for (var i=0; i<inputs.length; i++) {
    var val = parseFloat(inputs.eq(i).val());
    total+=val;
  }
  $("#weekly-total").text(total);
}

function hideEntryDialog() {
  $("#begin").val("");
  $("#end").val("");
  $("#entry_id").val("");
  $("#create_type").val("create");
  $("#add_entry_dialog").modal("toggle");
}
