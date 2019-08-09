package main

const APIPort = ":8080"
const DBName = "./hemato.db"

func main() {

	a := App{}
	a.Initialize(DBName)
	a.Run(APIPort)

}
