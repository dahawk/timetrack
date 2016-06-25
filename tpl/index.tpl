<html>
  <head>
      <title>TimeTrack</title>
        {{template "head" }}
        <script>
          $(document).ready(function(){
            setInterval(updateTime,1000);
            jQuery.datetimepicker.setLocale('de');

            $("#add_entry_submit").click(submitAddEntry);
            $("#refresh_entries").click(loadLogs);
            $("#active_entry").click(activeEntry);
            $("#add_entry").click(function() {
              $("#begin").val("");
              $("#end").val("");
            })
            $(".datepicker").datetimepicker({
              timepicker: false,
              format: 'd.m.Y'
            });

            var f = function(current) {
              if ($("#type").val() === 'Work time') {
                 this.setOptions({
                   timepicker: true,
                   format: 'd.m.Y H:i'
                 });
              } else {
                this.setOptions({
                  timepicker: false,
                  format: 'd.m.Y'
                });
              }
            }

            $(".datetimepicker").datetimepicker({
              format: 'd.m.Y H:i',
              onShow: f
            });

            loadLogs();
          });
        </script>
  </head>
  <body>
    {{template "navbar" .}}
    <div class="container">
      <div class="row timeselect">
        <div class="col-md-8 col-md-offset-2">
          <div class="row">
            <form class="form-inline" id="from_to_form">
              <div class="col-md-4 form-group"><label for="from_date">From</label><input type="text" name="from_date" id="from_date" class="datepicker" value="{{.From}}"></div>
              <div class="col-md-4 form-group"><label for="to_date">To</label><input type="text" name="to_date" id="to_date" class="datepicker" value="{{.To}}"></div>
              <div class="col-md-1"><button id="refresh_entries" class="btn btn-primary" type="button">Refresh</button></div>
            </form>
          </div>
        </div>
      </div>
      <div id="add_entry_dialog"  class="modal fade" role="dialog" data-keyboard="false">
         <div class="modal-dialog">
           <div class="modal-content">
             <div class="modal-header">
               <button type="button" class="close" onclick="hideEntryDialog">&times;</button>
               <h5>Add Time Entry</h5>
             </div>
             <div class="modal-body">
               <form role="form" id="add_entry_form">
                 <div class="form-group">
                   <label for="type">Type:</label>
                   <select class="form-control" id="type" name="type">
                     <option value="Work time">Work time</option>
                     <option value="Holiday">Holiday</option>
                     <option value="Sick leave">Sick leave</option>
                   </select>
                  </div>
                 <div class="form-group">
                  <label for="begin">Begin:</label>
                  <input type="text" class="form-control datetimepicker" id="begin" name="begin">
                </div>
                <div class="form-group">
                 <label for="end">End:</label>
                 <input type="text" class="form-control datetimepicker" id="end" name="end">
               </div>
               <input type="hidden" id="create_type" name="create_type" value="create"/>
               <input type="hidden" id="entry_id" name="entry_id"/>
               </form>
             </div>
             <div class="modal-footer">
               <button type="button" class="btn btn-primary" id="add_entry_submit">Save</button>
               <button type="button" class="btn btn-default" onclick="hideEntryDialog()">Close</button>
             </div>
         </div>
      </div>
    </div>
    <div class="container">
      <div class="row" id="stats-block">
        <div class="col-md-3 col-md-offset-3">
          <p><strong>Expected Work Time:</strong> {{.Stats.ExtectedWorkTime}}</p>
          <p><strong>Actual Work Time:</strong> {{.Stats.ActualWorkTime}}</p>
          <p><strong>Difference:</strong> {{.Stats.Delta}}</p>
        </div>
        <div class="col-md-3">
          <p><strong>Holidays:</strong> {{.Stats.Holidays}}</p>
          <p><strong>Sickdays:</strong> {{.Stats.Sickdays}}</p>
        </div>
      </div>
    </div>
    <table class="table table-condensed table-hover">
      <thead>
          <tr>
            <td>Action</td>
            <td colspan="2">Start</td>
            <td colspan="2">End</td>
            <td>Duration</td>
            <td>Type</td>
          </tr>
      </thead>
      <tbody>
      </tbody>
    </table>
  </div>
  </body>
</html>
