# TTK4145 Real-time programming - Elevator project
## Henrik Nyholm og Herman K. Dieset

![NTNU](https://innsida.ntnu.no/c/wiki/get_page_attachment?p_l_id=22780&nodeId=24647&title=Bruksregler+for+NTNU-logoen&fileName=variant1.jpg)

This is our elevator project, made in Golang with love. We hope you like it :) 
Features:
  - Very nice
  - Works really well
  - I've just learned to write README.md files.

```none
$$\   $$\                               $$\ $$\                                                               
$$ |  $$ |                              \__|$$ |                                                              
$$ |  $$ | $$$$$$\  $$$$$$$\   $$$$$$\  $$\ $$ |  $$\        $$$$$$$\ $$\   $$\  $$$$$$\   $$$$$$\   $$$$$$\  
$$$$$$$$ |$$  __$$\ $$  __$$\ $$  __$$\ $$ |$$ | $$  |      $$  _____|$$ |  $$ |$$  __$$\ $$  __$$\ $$  __$$\ 
$$  __$$ |$$$$$$$$ |$$ |  $$ |$$ |  \__|$$ |$$$$$$  /       \$$$$$$\  $$ |  $$ |$$ /  $$ |$$$$$$$$ |$$ |  \__|
$$ |  $$ |$$   ____|$$ |  $$ |$$ |      $$ |$$  _$$<         \____$$\ $$ |  $$ |$$ |  $$ |$$   ____|$$ |      
$$ |  $$ |\$$$$$$$\ $$ |  $$ |$$ |      $$ |$$ | \$$\       $$$$$$$  |\$$$$$$  |\$$$$$$$ |\$$$$$$$\ $$ |      
\__|  \__| \_______|\__|  \__|\__|      \__|\__|  \__|      \_______/  \______/  \____$$ | \_______|\__|      
                                                                                $$\   $$ |                    
                                                                                \$$$$$$  |                    
                                                                                 \______/                     
```

## Roadmap - Modules
- [ ] Main (to start goroutines)
- [x] Network (communicating to other elevators)
- [ ] Sync array (updates and modifies the sync array)
- [ ] Cost function (calculates whether an elevator should accept an order)
- [ ] IO module (talks with the physical elevator)

## Installation and useful sidenotes

Golang has to be installed, and the elevator's drivers (can be found through the TTK4145 repository).

To set the GOPATH environment value in Linux, the following lines of code has to be run:

```
$ cd
$ vim .bashrc                                           // Or other editor
$ export GOPATH=$HOME/gruppeOgPlass9/project-gruppe-9/  // Set Gopath to project folder that contains the src-folder. 
$ source .bashrc                                        // saves
$ go env                                                // displays the GOPATH
```
