{{range .}}<tr>
  <td>
    {{if .Active}}<button data-id="{{.EntryID}}" onclick="editEntry(this)" type="button"><i class="fa fa-pencil-square-o"></i></button>&nbsp;{{end}}
    <button data-id="{{.EntryID}}" onclick="deleteEntry(this)" type="button"><i class="fa fa-trash"></i></button>
  </td>
  <td>{{.DateFrom}}</td>
  <td>{{.TimeFrom}}</td>
  <td>{{.DateTo}}</td>
  <td>{{.TimeTo}}</td>
  <td>{{.Duration}}</td>
  <td>{{.Type}}</td>
</tr>{{end}}
