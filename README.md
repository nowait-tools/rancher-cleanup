## Rancher-cleanup

### Usage

Rancher-cleanup is available on docker hub [here](https://hub.docker.com/r/nowait/rancher-cleanup)

## Running on Rancher

Use following `docker-compose.yml`
```yml
rancher-cron:
  labels:
    io.rancher.container.create_agent: 'true'
    io.rancher.container.agent.role: environment
  image: socialengine/rancher-cron:0.2.0
```

It is important to include both labels as Rancher will set `CATTLE_URL`, 
`CATTLE_ACCESS_KEY`, and `CATTLE_SECRET_KEY`. If you want a bit more control,
feel free to set those manually.

The rancher-cleanup image sets sane defaults so as long as the labels mentioned above are set properly the image will run as is.  We recommend using [rancher-cron](https://github.com/SocialEngine/rancher-cron) to make this a scheduled job in order to keep the reconnecting hosts clean.
