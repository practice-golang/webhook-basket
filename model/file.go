package model

import "os"

var FileConnections, FileRequests *os.File

var TempClonedRepoRoot = "repos"         // Temporary Cloned repository root
var DeploymentRoot = "/home/my-websites" // Default deployment target root at remote server
