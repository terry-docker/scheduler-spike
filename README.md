# Spike for getting a go task scheduler

### http test commands

`$ curl -X POST http://localhost:8080/addTask     -H "Content-Type: application/json"     -d '{"Spec": "* * * * *", "VolumeName": "Every minute task"}'`

returns `{"taskID":1}`


`$ curl -X POST http://localhost:8080/removeTask     -H "Content-Type: application/json"     -d '{"TaskID": 1}'`


`curl -X GET http://localhost:8080/debug/tasks`

Mock gen:

`mockgen -source=persistence/persistence.go -destination=mocks/persistence/mock_persistence.go -package=mock_persistence`
`mockgen -source=cronjobs/cron.go -destination=mocks/cronjobs/mock_cron.go -package=mock_cronjobs`
