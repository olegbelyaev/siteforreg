package main

import (

    "github.com/gin-gonic/gin"
    "fmt"
    "os"
    "os/exec"
    "log"
)


func main() {
    secret:=os.Getenv("WEBHOOK_SECRET")
    if len(secret)==0 {
	panic("WEBHOOK_SECRET env not defined")
    }

    r := gin.Default()

    r.GET("/webhook/:action/:secret", func(c *gin.Context) {
	givenAction := c.Param("action")
	givenSecret := c.Param("secret")
	webhook(c,givenAction,givenSecret,secret)
    })

    r.POST("/webhook", func(c *gin.Context){
	givenAction := c.PostForm("action")
	givenSecret := c.PostForm("secret")
	webhook(c,givenAction,givenSecret,secret)
    } )
    
    r.Run(":5001") 
}

func webhook(c *gin.Context, givenAction string, givenSecret string, secret string){

	if len(givenAction)==0 {
	    c.JSON(400, gin.H{
	        "message": "empty action",
	    })
	    return
	}

	if givenSecret!=secret{
	    c.JSON(400, gin.H{
	        "message": "bad secret",
	    })
	    return
	}
	
	cmd := exec.Command("bash", "-c", fmt.Sprintf("./webhook.sh %s 2>&1", givenAction))
	
	stdoutStderr, err := cmd.CombinedOutput()
        
	c.JSON(500, gin.H{
	    "message": fmt.Sprintf("%s", stdoutStderr),
	})
	
	if err != nil {
	    log.Printf("%s", stdoutStderr)
	}

}