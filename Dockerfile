FROM alpine
WORKDIR /opt/swarmeus
COPY swarmeus_linux_amd64 /opt/swarmeus/
COPY swarmeus.yml /etc/swarmeus/swarmeus.yml
EXPOSE 9723
CMD ["/opt/swarmeus/swarmeus_linux_amd64", "--cfg.path=/etc/swarmeus/swarmeus.yml"]
