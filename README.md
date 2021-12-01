# Go Setup

## Setup go in Docker
``` bash
mkdir -p ~/bin
ln -s $PWD/goid.sh ~/bin/go
```

## Your first program
see https://go.dev/doc/code
``` bash
modulePath=<module-path> # Example: mby.fr/mass
mkdir -p src/$modulePath
cd src/$modulePath
go mod init $modulePath
```

## Using cobra
``` bash
go get github.com/spf13/cobra/cobra
go install github.com/spf13/cobra/cobra
ln -s $PWD/bin/cobra ~/bin/cobra
export GOPATH=$PWD
cobra
```

# Testing
To launch all unit tests there is a conveniant script:
``` bash
./unittest.sh
```

To launch unit test on an entire go module:
``` bash
cd <my_go_module>
go test -cover ./...
```


# Specifications

## Project management
### mass init workspace <path>
Initialize a new workspace.
- Create a new workspace directory
- Create a hidden cache dire .mass
- Create a settings file mass.yaml
- Create the default minimal config for dev, stage and prod envs

### mass init env <name>
Initialize a new environment.
- Create a new env directory

### mass init project <name>
Initialize a new project.
- Create a new project directory


## Project versionning

### mass version <project>
Display all versions in a project ?

### mass version bump <version>
Bump a version.


