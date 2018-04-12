package show

import "testing"

func TestParseFromHTML(t *testing.T) {
	var html = `<!DOCTYPE html>
<html>
<body>

<table style="width:100%">
  <tr>
    <th>Firstname</th>
    <th>Lastname</th>
    <th>Age</th>
    <td class="print_ignore"></td>
  </tr>
  <tr>
    <td>Jill</td>
    <td>Smith</td>
    <td>50</td>
  </tr>
  <tr>
    <td>Eve</td>
    <td>Jackson</td>
    <td>94</td>
  </tr>
  <tr>
    <td>John</td>
    <td>Doe</td>
    <td>80</td>
  </tr>
</table>

</body>
</html>`
	header, body, err := parseFromHTML(html)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v\n", header)
	t.Logf("%#v\n", body)
}
