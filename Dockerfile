FROM registry.cn-hangzhou.aliyuncs.com/wise2c/deploy-ui:v0.1

WORKDIR /deploy
VOLUME [ "/deploy" ]

COPY wise-deploy .
COPY table.sql .

ENTRYPOINT [ "bash", "-c", "/root/entrypoint.sh &\n ./wise-deploy -logtostderr" ]