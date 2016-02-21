<html>
  <head>
    <title>TimeTrack Login</title>
    {{template "head" }}
  </head>
  <body>
    {{template "navbar" .}}
    <div class="container">
      <h1>Login</h1>
      {{if .Error}}<div class="alert alert-danger">
        {{.Error}}
      </div>{{end}}
      <form class="form form-horizontal" role="form" method="POST" action="/login">
        <div class="form-group">
          <label for="user" class="control-label col-md-2">User</label>
          <div class="col-md-10">
            <input type="text" name="user" id="user"/>
          </div>
        </div>
        <div class="form-group">
          <label for="password" class="control-label col-md-2">Password</label>
          <div class="col-md-10">
            <input type="password" name="password" id="password"/>
          </div>
        </div>
        <button class="btn btn-primary" type="submit">Login</button>
      </form>
    </div>
  </body>
</html>
