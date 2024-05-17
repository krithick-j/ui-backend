package tmplts

var ICouponTemplate string = `
<!DOCTYPE html>
<html>
<head>
    <title>Your new iCoupon</title>
    <style>
        body {
            font-family: Arial, sans-serif;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }
        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
    <h2>This is a noreply email. Your iCoupons are:</h2>
    <table>
        <thead>
            <tr>
                <th>VID</th>
                <th>PIN</th>
                <th>Value</th>
                <th>Date On</th>
                <th>Expires On</th>
                <th>Active</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr>
                <td>{{.VID}}</td>
                <td>{{.Pin}}</td>
                <td>{{.Value}}</td>
                <td>{{.DateOn}}</td>
                <td>{{.ExpiresOn}}</td>
                <td>{{.Active}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`
