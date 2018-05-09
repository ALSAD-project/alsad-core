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
        label = float(lp[lp.find('(') + 1: lp.find(')')])
        vec = Vectors.dense(lp[lp.find('[') + 1: lp.find(']')].split(','))

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

    model = StreamingKMeans(k=5, decayFactor=1.0).setRandomCenters(4, 1.0, 0)

    model.trainOn(trainingData)

    result = model.predictOn(trainingData)
    result.pprint()

    result = model.predictOnValues(testingData.map(lambda lp: (lp.label, lp.features)))
    result.pprint()

    ssc.start()
    ssc.awaitTermination()