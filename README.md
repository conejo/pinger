# pinger

## description
using dhcp and running multiple pi's / bones / etc, and can't keep track of what ip address they're using? 

pinger reads the hostname and ip addresses of the machine during boot. Using the twilio api, it sends this information via SMS to the configured phone number. 

## usage

```.\pinger.exe -configFile settings.json```

pinger reads necessary config values from the supplied config filename param. To set these values, create a settings.json file using settings.sample.json as a guide.

set pinger to run at boot time via upstart or systemd. example .conf and .service files are provided.
