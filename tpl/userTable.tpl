{{range .}}
<tr>
  <td>
    <button data-id="{{.UserID}}" onclick="editUser(this)" type="button"><i class="fa fa-pencil-square-o"></i></button>&nbsp;
    {{if .Enabled}}
    <button data-id="{{.UserID}}" onclick="disableUser(this)" type="button"><i class="fa fa-lock"></i></button>
    {{else}}
    <button data-id="{{.UserID}}" onclick="enableUser(this)" type="button"><i class="fa fa-unlock"></i></button>
    {{end}}
  </td>
  <td>{{.Username}}</td>
  <td>{{.Name}}</td>
  <td>{{if .Username}}true{{else}}false{{end}}</td>
</tr>
{{end}}
