{{define "head"}}

<meta name="viewport" content="width=device-width, initial-scale=1">

<!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>

<!-- Latest compiled and minified CSS -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">

<!-- Optional theme -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css" integrity="sha384-fLW2N01lMqjakBkx3l/M9EahuwpSfeNvV63J5ezn3uZzapT0u7EYsXMjQV+0En5r" crossorigin="anonymous">

<link rel="stylesheet" type="text/css" href="/static/datetimepicker/jquery.datetimepicker.css"/>
<script src="/static/datetimepicker/build/jquery.datetimepicker.full.min.js"></script>
<script src="/static/timetracker.js"></script>
<link rel="stylesheet" type="text/css" href="/static/timetrack.css"/>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.5.0/css/font-awesome.min.css"/>


<!-- Latest compiled and minified JavaScript -->
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js" integrity="sha384-0mSbJDEHialfmuBBQP6A4Qrprq5OVfW37PRR3j5ELqxss1yVqOtnepnHVP9aJ7xS" crossorigin="anonymous"></script>

{{end}}
{{define "navbar"}}<nav class="navbar navbar-default">
  <div class="container-fluid">
    <div class="navbar-header">
      <a class="navbar-brand" href="/">TimeTrack</a>
    </div>
    <ul class="nav navbar-nav">
      {{if .User.UserID}}
      <li>
        <div class="btn-group">
          <button type="button" class="btn btn-success navbar-padding" id="active_entry" data-active="{{.User.ActiveEntry}}">{{if .User.ActiveEntry}}Stop{{else}}Start{{end}}</button>
        </div>
      </li>
      <li>&nbsp;</li>
      <li>
        <div class="btn-group">
          <button type="button" class="btn btn-default navbar-padding" id="add_entry" data-toggle="modal" data-target="#add_entry_dialog">Add Entry</button>
          <a href="/pdf?dateFrom={{.From}}&dateTo={{.To}}" class="btn btn-default navbar-padding" target="_blank" id="pdf-btn">PDF</a>
        </div>
      </li>
      {{end}}
    </ul>
    <ul class="nav navbar-nav navbar-right">
      <li><div class="navbar-padding"><span id="time"></span></div></li>
      <li class="dropdown">
        <a class="dropdown-toggle" data-toggle="dropdown" href="#"><span class="glyphicon glyphicon-user"></span> {{.User.Username}}
        <span class="caret"></span></a>
        <ul class="dropdown-menu">
          {{if .User.Admin}}<li><a href="/admin">Admin</a></li>{{end}}
          {{if .Impersonating}}<li><a href="/unimpersonate">Unimpersonate</a></li>{{end}}
          <li><a href="/logout">Logout</a></li>
        </ul>
      </li>
    </ul>
  </div>
</nav>{{end}}

{{define "navbar-admin"}}<nav class="navbar navbar-default">
  <div class="container-fluid">
    <div class="navbar-header">
      <a class="navbar-brand" href="/admin">TimeTrack Admin</a>
    </div>
    <ul class="nav navbar-nav">
      {{if .User.UserID}}
      <li>
        <div class="btn-group">
          <a href="/admin" class="btn btn-default navbar-padding">Users</a>
        </div>
      </li>
      {{end}}
    </ul>
    <ul class="nav navbar-nav navbar-right">
      <li><div class="navbar-padding"><span id="time"></span></div></li>
      <li class="dropdown">
        <a class="dropdown-toggle" data-toggle="dropdown" href="#"><span class="glyphicon glyphicon-user"></span> {{.User.Username}}
        <span class="caret"></span></a>
        <ul class="dropdown-menu">
          <li><a href="/logout">Logout</a></li>
        </ul>
      </li>
    </ul>
  </div>
</nav>{{end}}

{{define "stats"}}
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
{{end}}
