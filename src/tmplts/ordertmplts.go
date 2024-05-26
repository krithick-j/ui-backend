package tmplts

var OrderTemplate string = `<!DOCTYPE html>
<html>
<head>
    <title>Your Order Detail</title>
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
    <p>This is a noreply email. Your Order Details are:</p>
	<p>Dear {{.CustomerDetails.Name}}, </p>
	<p>We received your order on <date>. You have 15 days refund window from the date of order</p>
	<p>Your order details as below</p>
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Qty</th>
                <th>Unit Price</th>
                <th>Total</th>
            </tr>
        </thead>
        <tbody>
            {{range .ProductDetails}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.Quantity}}</td>
                <td>{{.Price}}</td>
                <td>{{.SubTotal}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`
