# flightpath

This is a trivial attempt to find an iternary from a given list of aiport codes

The API /calculate returns an Iternary from a given list of inputs. 


```
curl -XPOST localhost:8080/calculate -d@/tmp/input.json

$ cat /tmp/input.json 
{
  "A": "B",
  "B": "C",
  "C": "D",
  "D": "E"
}

```

Caveats
1. A simple attempt is made to detect cycles or incorrect inputs.
2. The logic does not cover the usecase of disjointed paths.

TODO:
Create a docker image and create a simple K8s Service that can listen on 8080.
Due to lack of time, did not do this yet.

