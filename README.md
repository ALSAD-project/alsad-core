ALSAD Core
===========

This is the core engine for ALSAD.

## Project Structure

```
.
├── README.md              # readme file of the project
├── cmd                    # executable entry points
│   ├── datafeeder             # data feeder
│   ├── expertsystem           # expert system
│   │   ├── daemon                 # daemon of the expert system  
│   │   ├── terminal               # expert system terminal
│   │   └── web                    # expert system web interface
│   └── learners           # base learner modules
│       └── drivers            # learner drivers
│           └── spark              # learner driver for Apache Spark
├── docker-compose.yml     # docker-compose file
├── dockerfiles            # docker files for building images
│   └── spark.dockerfile       # a minimal Apache Spark on Mesos
├── pkg                    # packages
│   ├── datafeeder             # data feeder package
│   ├── expertsystem           # expert system package
│   └── learners               # learner package
│       └── drivers                # learner driver package
│           └── spark                  # learner driver package for Apache Spark    
└── scripts                # some useful scripts

```
