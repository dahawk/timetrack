<html>
  <head>
    <title>TimeTrack</title>
    {{template "head" }}
    <script>
      $(document).ready(function() {
        $("#fulltime").click(function(){
          if ($(this).prop("checked")) {
            $("#parttime").hide();
            $(".day_input").attr("value","7.7");
          } else {
            updateTotal();
            $("#parttime").show();
          }
        });

        $(".day-input").change(updateTotal);
      });
    </script>
  </head>
  <body>
    {{template "navbar-admin" .}}
    <div class="container">
      <h1>{{.WorkTimeUser.Name}}'s Work Time</h1>
      <br/>
      <form action="/worktime" method="POST">
      <div class="row">
        <div class="col-md-2 col-md-offset-1">
          <label for="fulltime">Full time</label>
          <input type="checkbox" name="fulltime" id="fulltime" {{if .WorkTime.FullTime}}checked{{end}} value="fulltime"/>
        </div>
      </div>
      <div class="row{{if .WorkTime.FullTime}} collapse{{end}}" id="parttime">
        <table class="table table-condensed">
          <thead>
            <tr>
              <td>Mon</td>
              <td>Tue</td>
              <td>Wed</td>
              <td>Thu</td>
              <td>Fri</td>
              <td></td>
              <td><strong>Weekly Total</strong></td>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><input type="text" id="mon-input" class="day-input" name="mon-input" value="{{.WorkTime.Mon}}"/></td>
              <td><input type="text" id="tue-input" class="day-input" name="tue-input" value="{{.WorkTime.Tue}}"/></td>
              <td><input type="text" id="wed-input" class="day-input" name="wed-input" value="{{.WorkTime.Wed}}"/></td>
              <td><input type="text" id="thu-input" class="day-input" name="thu-input" value="{{.WorkTime.Thu}}"/></td>
              <td><input type="text" id="fri-input" class="day-input" name="fri-input" value="{{.WorkTime.Fri}}"/></td>
              <td></td>
              <td><strong><span id="weekly-total">{{.WorkTime.WeeklyTime}}</span></strong></td>
            </tr>
          </tbody>
        </table>
        <input type="hidden" id="user" name="user" value="{{.WorkTimeUser.UserID}}"/>
      </div>
      <div class="row">
        <div class="col-md-1">
          <button type="submit" class="btn btn-primary">Submit</button>
        </div>
      </div>
    </form>
  </body>
</html>
