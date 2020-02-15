# DigitalOcean Firewall Updater

With all credit to, extended and built on the work of

```
github.com/paolobarbolini/do-firewall-updater
```

## How it works

Update all your DigitalOcean firewall rules based on a json file. First run with no JSON files present will create the old and new files in the run folder.

Thereafter simply edit the downloaded rules in the ``new_ips.json`` file and if these are different to ``old_ips.json`` then the next run will update your firewall.

## Running it
Build the executable and run it as below

```
/path/to/executable --token DIGITALOCEAN_API_TOKEN --firewall-id THE_FIREWALL_ID
```

or

```
/path/to/executable --token DIGITALOCEAN_API_TOKEN --firewall-name THE_FIREWALL_NAME
```

To generate a new api token go to the [Applications & API section](https://cloud.digitalocean.com/settings/api/tokens) in the digitalocean control panel and create a new personal access token.
The token must have read and write privileges.

## Running it regularly
Create a cron job with:

```
crontab -e
*/15 * * * * /path/to/executable --token DIGITALOCEAN_API_TOKEN --firewall-id THE_FIREWALL_ID
```

This example will run every 15 minutes.

If you want to update the same firewall from multiple computers with different connections make sure that the cron job doesn't run at the same time as the other computers to prevent race conditions.
