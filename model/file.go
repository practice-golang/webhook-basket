package model

import "os"

var FileConnections, FileRequests *os.File

var CloneRepoRoot = "repos"              // Cloned repository root
var DeploymentRoot = "/home/my-websites" // Default deployment target root at remote server
