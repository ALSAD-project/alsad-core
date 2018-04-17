#!/usr/bin/env python

import pandas as pd
import numpy as np
import socket
from time import sleep
from sklearn import linear_model
from sklearn import preprocessing
from sklearn.model_selection import train_test_split

# Global settings
pd.options.display.max_columns = 999
pd.options.display.max_rows = 999

##### Begin of Initial Training #####

feature_col = ["sepal length", "sepal width", "petal length", "petal width", "class"]
predict_col = "class"
preprocess_col = ["class"]

df = []
df.append(pd.read_csv('data/shuffled_iris.csv', sep=',', header=None))
df = pd.concat(df).reset_index(drop=True)
df.columns = feature_col
df[preprocess_col] = df[preprocess_col].apply(preprocessing.LabelEncoder().fit_transform)

y = df.pop(predict_col)
X = df
X_train, X_incr, y_train, y_incr = train_test_split(X, y, test_size=0.5, shuffle=False)

clf = linear_model.SGDClassifier(max_iter=1000)
clf.fit(X_train, y_train)

##### End of Initial Training #####

##### Begin of Incremental Training #####

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect(("localhost", 8888))

while True:
	response = client.recv(4096)
	client.send(response)

##### End of Incremental Training #####

