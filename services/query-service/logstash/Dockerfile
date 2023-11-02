FROM docker.elastic.co/logstash/logstash:7.3.2

# install dependency
RUN /usr/share/logstash/bin/logstash-plugin install logstash-input-jdbc
RUN /usr/share/logstash/bin/logstash-plugin install logstash-filter-aggregate
RUN /usr/share/logstash/bin/logstash-plugin install logstash-filter-jdbc_streaming
RUN /usr/share/logstash/bin/logstash-plugin install logstash-filter-mutate

# copy lib database jdbc jars
COPY ./mysql-connector-java-8.0.15.jar /usr/share/logstash/logstash-core/lib/jars/mysql-connector-java.jar