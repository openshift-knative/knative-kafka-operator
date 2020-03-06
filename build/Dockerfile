FROM registry-proxy.engineering.redhat.com/rh-osbs/ubi8-minimal:8-released

ENV OPERATOR=/usr/local/bin/knative-kafka-operator \
    USER_UID=1001 \
    USER_NAME=knative-kafka-operator
    
# install operator binary
COPY build/_output/bin/knative-kafka-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

# install manifest[s]
COPY deploy /deploy

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
