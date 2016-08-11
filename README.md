# Open Sesame for Dimension Data Cloud Control

Poke a hole in the firewall for a given network domain from your current location.

## Usage

Configure your credentials (on Windows, replace `export` with `set`):
```bash
export DD_COMPUTE_USER="my_user_name"
export DD_COMPUTE_PASSWORD="my_password"
```

Specify the target network domain by its name:
```
dd-sesame --region AU --network Mantl2 --dc AU9 --rule me.at.home.ipv4.inbound
```

Specify the target network domain by its Id:
```
dd-sesame --region AU --network d35e1cf5-a8cd-48f7-b4d6-075d51fde461 --rule me.at.home.ipv4.inbound
```

Delete the firewall rule once you're done:
```
dd-sesame --delete --region AU --network Mantl2 --dc AU9 --rule me.at.home.ipv4.inbound
```
