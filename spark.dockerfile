FROM mesosphere/spark:2.1.0-2.2.0-1-hadoop-2.7

ENV SPARK_HOME /opt/spark/dist
ENV PATH $SPARK_HOME/bin:$PATH
