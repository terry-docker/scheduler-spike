# Spike for getting a go task scheduler

### http test commands

`$ curl -X POST http://localhost:8080/addTask     -H "Content-Type: application/json"     -d '{"Spec": "* * * * *", "Description": "Every minute task"}'`

returns `{"taskID":1}`


`$ curl -X POST http://localhost:8080/removeTask     -H "Content-Type: application/json"     -d '{"TaskID": 1}'`