package tmplts

var DistribApplicationFormTemplate string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Distributor Application Form</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        h1, h2 {
            text-align: center;
        }
        .section {
            margin-bottom: 20px;
        }
        .section-title {
            font-weight: bold;
            text-decoration: underline;
        }
        .field-label {
            display: inline-block;
            width: 200px;
            font-weight: bold;
        }
        .field-value {
            display: inline-block;
        }
    </style>
</head>
<body>
    <h1>DISTRIBUTOR APPLICATION FORM</h1>
    
   
    
    <div class="section">
        <h2>APPLICATION INFORMATION</h2>
        <div>
            <div class="field-label">Title (If Individual):</div>
            <div class="field-value">{{.Title}}</div>
        </div>
        <div>
            <div class="field-label">Surname, Given Name (If Individual):</div>
            <div class="field-value">{{.Name}}</div>
        </div>
        <div>
            <div class="field-label">Mailing Address:</div>
            <div class="field-value">{{.MailingAddress}}</div>
        </div>
        <div>
            <div class="field-label">Cheque Name:</div>
            <div class="field-value">{{.ChequeName}}</div>
        </div>
        <div>
            <div class="field-label">Shipping Address:</div>
            <div class="field-value">{{.Address1}}</div>
        </div>
        <div>
            <div class="field-label">Home Phone No</div>
            <div class="field-value">{{.HomePhoneNo}}</div>
        </div>
        <div>
            <div class="field-label">Valid ID Type / ID Number:</div>
            <div class="field-value">{{.ValidIdNo}}</div>
        </div>
        <div>
            <div class="field-label">eMail Address:</div>
            <div class="field-value">{{.EmailAddress}}</div>
        </div>
        <div>
            <div class="field-label">Nationality & Date of Birth:</div>
            <div class="field-value">{{.Country}}; {{.DateOfBirth}}</div>
        </div>
        <div>
            <div class="field-label">Mother's Maiden Name:</div>
            <div class="field-value">{{.MothersMaidenName}}</div>
        </div>
        <div>
            <div class="field-label">Name of Beneficiary/Nominee:</div>
            <div class="field-value">{{.BenificiaryName}}</div>
        </div>
        <div>
            <div class="field-label">Relationship to Beneficiary/Nominee:</div>
            <div class="field-value">{{.BeneficiaryRelationship}}</div>
        </div>
        <div>
            <div class="field-label">PAN Card:</div>
            <div class="field-value">{{.PanCard}}</div>
        </div>
        <div>
            <div class="field-label">Bank Name:</div>
            <div class="field-value">{{.BankName}}</div>
        </div>
        <div>
            <div class="field-label">Bank Acct. No.:</div>
            <div class="field-value">{{.BankAccNo}}</div>
        </div>
        <div>
            <div class="field-label">IFS Code:</div>
            <div class="field-value">{{.IFSCCode}}</div>
        </div>
    </div>
    
    <div class="section">
        <h2>PREFERRED PLACEMENT INFORMATION</h2>
        <div>
            <div class="field-label">Direct Preferred Placement Distributor ID No:</div>
            <div class="field-value">{{.PreferredDistribId}}</div>
        </div>
        <div>
            <div class="field-label">Full Name:</div>
            <div class="field-value">{{.PreferredDistribName}}</div>
        </div>
        <div>
            <div class="field-label">Preferred Placement:</div>
            <div class="field-value">{{.PreferredPlace}}</div>
        </div>
        <div>
            <div class="field-label">Preferred Side:</div>
            <div class="field-value">{{.PreferredSide}}</div>
        </div>
    </div>
</body>
</html>
`

//  <div class="section">
//         <div class="field-label">Distributor ID No.:</div>
//         <div class="field-value">{{.DistribID}}</div>
//     </div>
//     <div class="section">
//         <div class="field-label">Date:</div>
//         <div class="field-value">{{.CreatedAt}}</div>
//     </div>

//     <div class="section">
//         <h2>REFERRER INFORMATION</h2>
//         <div>
//             <div class="field-label">Referrer ID:</div>
//             <div class="field-value">{{.RefDistribID}}</div>
//         </div>
//         <div>
//             <div class="field-label">Referrer Name:</div>
//             <div class="field-value">{{.RefDistribName}}</div>
//         </div>
//     </div>
