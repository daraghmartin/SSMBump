# SSMBump
Bump a semantic versioning variable in an SSM Parameter Store

./SSMBump "/Project/App/Version"

If the param does not exist you get 0.0.0.  If the param exists and is x.y.z then z will be incremented. 