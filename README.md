# ASGMaster

Jacob Verhoeks / jjverhoeks@edrans.com
October 2018 

Small go utility to have a command running only on one server in a ASG

1. Find my current ASG Role Tag
2. Check if there is any instance with asgrole:master tag
3. If there is none
  - Set tag to current instance
  - Return 0  
4. If the tag is our instanceID
  - Return 0  
5. If set to other machine
  - Return 1

## Commandline options  
```
  Usage of asgmaster:
    -debug
      	Enable Debugging
    -key string
      	Tag name to match the asg (default "role")
    -region string
      	AWS Region (default "eu-west-1")
```

## Prerequisites

The instance role need permissions to retreive and set tags

## usage

Cron:

`* * * * * root asgmaster && myscript.sh`
This will only execute if current node is master
