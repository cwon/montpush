package main

import (
	webpush "github.com/SherClockHolmes/webpush-go"
	"fmt"
)

func main() {
	privateKey, publicKey, _ := webpush.GenerateVAPIDKeys()
	
	fmt.Println("private : " + privateKey);
	fmt.Println("public  : " + publicKey);
}
