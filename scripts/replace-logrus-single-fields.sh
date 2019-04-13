#! /bin/bash

# Run this in the docker container
find . -name '*.go' | xargs sed -r -i 's|WithFields\(logrus\.Fields\{"([^"]+)": ([^},]+)}\)|WithField("\1", \2)|g'