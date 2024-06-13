task("build") {
    doLast {
        //GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
        exec {
            executable = "go"
            environment = environment
                .plus(Pair("GOOS", "linux"))
                .plus(Pair("GOARCH", "arm64"))
            args = listOf("build", "-o", "bootstrap", "./pkg/main.go")
        }
        exec {
            executable = "zip"
            args = listOf("lambdaFunction.zip", "bootstrap")
        }
    }
}
task<Delete>("clean") {
    delete("bootstrap")
    delete("lambdaFunction.zip")
}

task("test") {
    exec {
        executable = "go"
        args = listOf("test", "-v", "./...")
    }
}