Shows differeces between csv exports from a bank (currently verity) and
Quickbooks.

Usage:
```
go run account-diff.go [verity|amex] [export from bank] [export from quickbooks]
go run account-diff.go verity ~/Downloads/History.csv ~/Downloads/Register.csv
```

`History.csv` is from Verity
`Register.csv` is from Quickbooks


in QB:
Transactions -> Chart of Accounts -> account -> export
save exported file as csv

At bank
Export -> csv -> Year to date

amex
Export
split dates into two columns