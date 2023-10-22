#FROM docker.io/cytopia/ansible:latest-infra-0.34
FROM docker.io/cytopia/ansible:2.10-tools

RUN ansible-galaxy collection install containers.podman

WORKDIR /deploy
VOLUME [ "/deploy" ]

COPY pagoda /root
COPY database/table.sql /root
COPY favicon.ico /root

ENTRYPOINT [ "/root/pagoda", "-logtostderr", "-v", "4", "-w", "/workspace" ]
