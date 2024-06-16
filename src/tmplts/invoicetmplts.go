package tmplts

var InvoiceTemplate string = `<!DOCTYPE html>
<html>
    <head>
    </head>
    <body>
		<p>Dear {{.CustomerDetails.Name}}({{.DistribId}}), </p>
	    <p>We received your order on <date>. You have 15 days refund window from the date of order</p>
	    <p>Your Invoice details as below</p>
        <div id="print-area">
        <table width="800">
            <tr>
                <th colspan="6">
                <img style="float:left;background-color: #c3c3c3" src="https://ui-network.com/whiteuiLogo.svg" />
            <div style="font-size:1.8em;color:white; background-color:  #c3c3c3;text-align: center;padding-top:20px;height:65">Universe Network India</div>
                </th>
            </tr>
            <tr>
                <td colspan="6" align="center" style="font-weight: bold;font-size: 1.5em;">Invoice</td>
            </tr>
            <tr>
                <td colspan="3">To:<br />{{.CustomerDetails.Name}}<br />{{.ShippingAddress.Address}}<br />{{.ShippingAddress.City}}<br />{{.ShippingAddress.District}}<br />{{.ShippingAddress.State}}<br />{{.ShippingAddress.ZipCode}}<br />{{.ShippingAddress.Country}}</td>
                <td colspan="3">Ordered On:{{.CreatedAt}}<br />Invoice No: {{.OrderId}}<br />GST No. xxxxx yyyyy<br /></td>
            </tr>
            <tr>
                <th>Sl No</th>
                <th>Description</th>
                <th>Qty</td>
                <th>U. Price</th>
                <th>GST %</th>
                <th>Amount</th>
            </tr>
			{{range $index, $product := .Products}}
			<tr>
				<td align="middle">{{$index | incIndex}}</td>
				<td>{{$product.Name}}</td>
				<td align="middle">{{$product.Quantity}}</td>
				<td align="right">{{$product.UnitPrice}}</td>
				<td align="right">{{$product.GstPercentage}}</td>
				<td align="right">{{$product.SubTotal}}</td>
			</tr>
			{{end}}
			<tr style="font-weight: bold;">
                <td align="middle" colspan="5" >Total Shipping Cost</td>
                <td align="right">{{.TotalSandH}}</td>
            </tr>
            <tr style="font-weight: bold;">
                <td align="middle" colspan="5" >Grant Total</td>
                <td align="right">{{.TotalAmount}}</td>
            </tr>
        </table>
        </div>
        <button onclick="window.print()">Print</button>
    </body>
    <style>
        td {padding: 5px;}
        table, tr, th, td {border: 1px solid;border-collapse: collapse;}
        @media print {
  body {
    visibility: hidden;
  }
  #print-area {
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
    visibility: visible;
    position: absolute;
    left: 0;
    top: 0;
  }
}
    </style>
</html>
`
