# Usage: python3 buildauthheader.py <clientId> <clientSecret>

import sys
import urllib.parse
import base64
import os

client_id = urllib.parse.quote(sys.argv[1])
client_secret = urllib.parse.quote(sys.argv[2])

header = client_id + ":" + client_secret
header_bytes = header.encode('ascii')
header_base64 = base64.b64encode(header_bytes)
header = "Basic " + header_base64.decode('ascii')

print(header)
