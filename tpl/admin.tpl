<html>
  <head>
    <title>TimeTrack</title>
    {{template "head" }}
    <script>
      $(document).ready(function() {
        $("#edit_user_submit").click(submitUser);
        $("#add_user_btn").click(addUser);
      });
    </script>
  </head>
  <body>
    {{template "navbar-admin" .}}
    <div class="container">
      <h1>Users</h1>
      <div class="pull-right"><button id="add_user_btn" class="btn btn-default" type="button">Create User</button></div>
      <div id="edit_user_dialog"  class="modal fade" role="dialog">
         <div class="modal-dialog">
           <div class="modal-content">
             <div class="modal-header">
               <button type="button" class="close" data-dismiss="modal">&times;</button>
               <h5>Edit User</h5>
             </div>
             <div class="modal-body">
               <div class="alert alert-danger collapse" id="password-alert">
                 <a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a>
                  passwords don't match
               </div>
               <form role="form" id="user_edit_form">
                 <div class="form-group">
                   <label for="username">Username:</label>
                   <input type="text" class="form-control" id="username" name="username" />
                 </div>
                 <div class="form-group">
                   <label for="name">Name:</label>
                   <input type="text" class="form-control" id="name" name="name" />
                 </div>
                 <div class="form-group has-feedback" id="pwd_group">
                   <label for="password">Password:</label>
                   <input type="password" class="form-control" id="password" name="password" />
                 </div>
                 <div class="form-group has-feedback" id="rep_group">
                  <label for="repeat">Repeat:</label>
                  <input type="password" class="form-control" id="repeat" name="repeat" />
                </div>
                <input type="hidden" id="user_id" name="user_id"/>
                <input type="hidden" id="type" name="type"/>
               </form>
             </div>
             <div class="modal-footer">
               <button type="button" class="btn btn-primary" id="edit_user_submit">Save</button>
               <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
             </div>
         </div>
      </div>
    </div>
      <table class="table table-condensed table-hover">
        <thead>
          <tr>
            <td>Action</td>
            <td>Username</td>
            <td>Name</td>
            <td>Admin</td>
          </tr>
        </thead>
        <tbody>
          {{range .UserList}}
          <tr>
            <td>
              <div class="btn-group">
                <button data-id="{{.UserID}}" onclick="editUser(this)" type="button" class="btn btn-default"><i class="fa fa-pencil-square-o"></i></button>
                {{if .Enabled}}
                <button data-id="{{.UserID}}" onclick="disableUser(this)" type="button" class="btn btn-default"><i class="fa fa-lock"></i></button>
                {{else}}
                <button data-id="{{.UserID}}" onclick="enableUser(this)" type="button" class="btn btn-default"><i class="fa fa-unlock"></i></button>
                {{end}}
                <a href="worktime?id={{.UserID}}" class="btn btn-default"><i class="fa fa-clock-o"></i></a>
              </div>
            </td>
            <td><a href="impersonate?id={{.UserID}}" target="_blank">{{.Username}}</a></td>
            <td><a href="impersonate?id={{.UserID}}" target="_blank">{{.Name}}</a></td>
            <td>{{if .Username}}true{{else}}false{{end}}</td>
          </tr>
          {{end}}
        </tbody>
      </table>
    </div>
  </body>
</html>
