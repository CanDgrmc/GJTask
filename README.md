# GJTask

## Create 
http://localhost:3000/maze | POST
```json
{
	"arr": [
[1,1,1,1,0,1,1,1],
[1,0,0,0,0,0,0,1],
[1,0,1,1,1,0,1,1],
[1,0,0,0,1,0,0,1],
[1,1,1,0,1,1,0,1],
[1,0,0,0,1,0,0,1],
[1,0,1,1,1,0,1,1],
[1,0,0,2,0,0,0,1],
[1,1,1,1,1,1,1,1]
]

}
```

## GET
http://localhost:3000/maze/{id} | GET

## Calculate
http://localhost:3000/maze/{id}/calculate | GET

## DELETE
http://localhost:3000/maze/{id} | DELETE

