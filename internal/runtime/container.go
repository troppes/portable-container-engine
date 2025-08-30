package runtime

type ContainerRuntime interface {
    Run(image string, command []string) error
    CreateChildProcess(path string, command []string) error
}

func GetRuntime() ContainerRuntime {
    return &platformRuntime{}
}
