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
    webhookScript:=os.Getenv("WEBHOOK_SCRIPT")
    
    if len(secret)==0 {
	panic("WEBHOOK_SECRET env not defined")
    }

    r := gin.Default()

    r.GET("/webhook/:action/:secret", func(c *gin.Context) {
	givenAction := c.Param("action")
	givenSecret := c.Param("secret")
	webhook(c,givenAction,givenSecret,secret,webhookScript)
    })

    r.POST("/webhook", func(c *gin.Context){
	givenAction := c.PostForm("action")
	givenSecret := c.PostForm("secret")
	webhook(c,givenAction,givenSecret,secret,webhookScript)
    } )
    
    r.Run(":5001") 
}

func webhook(c *gin.Context, givenAction string, givenSecret string, secret string, webhookScript string){

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
	
	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s %s 2>&1",webhookScript, givenAction))
	
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
	    log.Printf("%s", stdoutStderr)
	}

        
	c.JSON(200, gin.H{
	    "message": fmt.Sprintf("%s", stdoutStderr),
	})
	

}