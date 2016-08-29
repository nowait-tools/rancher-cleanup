FROM golang:1.7.0-alpine

ENV RANCHER_API_TIMEOUT=10 \
    HOST_REMOVAL_INTERVAL=30
    
COPY rancher-cleanup /rancher-cleanup
RUN chmod +x /rancher-cleanup

ENTRYPOINT ["/rancher-cleanup"]
