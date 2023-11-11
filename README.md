# Fast VPN Checker

## Overview

Fast VPN Checker is a Go program that helps you find fast VPN servers by fetching a list of VPN servers from [VPN Gate](https://www.vpngate.net/en/) and measuring their response times through pinging. It's especially useful when the CSV endpoint (`CsvEndPoint`) is blocked in certain countries, as users can connect to VPN to access it.

## Prerequisites

- Go installed on your machine.
- Internet connectivity to fetch VPN server data.
- Access to external servers for pinging.

## Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/abdoroot/fast-vpn-checker.git
   cd fast-vpn-checker
2. Run the program:
   go run main.go
3. View the sorted list of fast VPN servers in the program's output.


Important Note

    The CSV endpoint (CsvEndPoint) may be blocked in some countries. If you encounter issues fetching data, consider connecting to a VPN before running the program.

Additional Information
 VPN Gate Website: https://www.vpngate.net/en/
