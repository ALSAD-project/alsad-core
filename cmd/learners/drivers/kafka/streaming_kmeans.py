# Example use:
# $ kubectl port-forward `kubectl get pods -o=jsonpath='{range .items[*]}{.metadata.name}{"\n"}' | grep driver` 4040:4040
# $ kubectl exec -ti kafka-0 -- kafka-console-producer.sh --topic streamin --broker-list localhost:9093
# $ kubectl exec -ti kafka-0 -- kafka-console-producer.sh --topic streamout --broker-list localhost:9093

from __future__ import print_function

import sys

from pyspark import SparkContext
from pyspark.streaming import StreamingContext
from pyspark.streaming.kafka import KafkaUtils

from pyspark.mllib.linalg import Vectors
from pyspark.mllib.regression import LabeledPoint
from pyspark.mllib.clustering import StreamingKMeans

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: streaming_kmeans.py <zk> <trainTopic> <testTopic>", file=sys.stderr)
        exit(-1)

    def parse(lp):
        record = [float(x) for x in lp.strip().split(',')]
        label = record[-1]
        vec = Vectors.dense(record[:-1])

        return LabeledPoint(label, vec)

    sc = SparkContext(appName="StreamingKmeansKafka")
    ssc = StreamingContext(sc, 20)

    zkQuorum, trainTopic, testTopic = sys.argv[1:]
    kvs = KafkaUtils.createStream(ssc, zkQuorum, "spark-streaming-consumer", {trainTopic: 1})
    lines = kvs.map(lambda x: x[1])
    trainingData = lines.map(lambda line: Vectors.dense([float(x) for x in line.strip().split(',')]))
    
    kvs = KafkaUtils.createStream(ssc, zkQuorum, "spark-streaming-consumer", {testTopic: 1})
    lines = kvs.map(lambda x: x[1])
    testingData = lines.map(parse)

    model = StreamingKMeans(k=3, decayFactor=1.0).setRandomCenters(4, 1.0, 0)

    model.trainOn(trainingData)

    result = model.predictOn(trainingData)
    result.pprint()

    result = model.predictOnValues(testingData.map(lambda lp: (lp.label, lp.features)))
    result.pprint()

    ssc.start()
    ssc.awaitTermination()