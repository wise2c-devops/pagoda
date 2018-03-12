FROM generik/ansible:v2.4

WORKDIR /deploy
VOLUME [ "/deploy" ]

COPY wise-deploy /root
COPY database/table.sql /root
COPY favicon.ico /root

ENTRYPOINT [ "/root/pagoda", "-logtostderr", "-v", "4", "-w", "/workspace" ]