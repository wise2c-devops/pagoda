FROM registry.cn-hangzhou.aliyuncs.com/wise2c/deploy-ui:v0.1

WORKDIR /deploy
VOLUME [ "/deploy" ]

COPY wise-deploy /root
COPY table.sql /root

ENTRYPOINT [ "bash", "-c", "/root/entrypoint.sh &\n /root/wise-deploy -logtostderr -v 4 -w /workspace" ]