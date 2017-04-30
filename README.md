# pinger

## description
pinger reads the hostname and ip addresses of the machine. Using the twilio api, it sends this information via SMS to the configured phone number. 

## usage
pinger reads necessary config values from the supplied config filename param. To set these values, create a settings.json file using settings.sample.json as a guide.

set pinger to run at boot time via cron or similar, it will then sms hostname and ip addresses on boot.
