FROM generik/ansible:2.3

COPY wise-deploy .

ENTRYPOINT [ "./wise-deploy", "-logtostderr" ]