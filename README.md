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

## Development
### Expert System
1. Please enable the `localConfigure()` in `pkg/expertsystem/expertsystem.go` for localhost setting
2. `$ go run cmd/expertsystem/daemon/main.go`
3. `$ go run cmd/expertsystem/terminal/main.go`

For running in docker:
1. `$ docker-compose up --build expertsystem-daemon`
2. `$ docker-compose up --build expertsystem-terminal`