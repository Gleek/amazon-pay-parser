# Amazon Pay Transaction parser

This is a simple parser for Amazon Pay transactions. It reads a CSV file with the transactions and outputs a CSV file with the parsed transactions.

Amazon pay transactions only seem to be present in their web interface while I needed them to be in an excel file to track my expenses better

## Usage
- go install github.com/gleek/amazon-pay-parser
- Extract the raw html from the Amazon Pay page.
  - Scroll to get enough transactions and inspect element to get div with id `transactions-desktop`
  - Copy the inner html of the div and save it to a file
- Run the parser
- `amazon-pay-parser transactions.html > transactions.csv`
