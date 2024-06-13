extra["yarnVersion"] = "1.22.19"
extra["nodeVersion"] = "20.10.0"

plugins {
    val nodePluginVersion = "5.0.0"
    id("com.github.node-gradle.node") version nodePluginVersion apply false
}

tasks.register("build") {
    dependsOn("lambda:build")
    dependsOn("cdk:build")
}

tasks.register("clean") {
    dependsOn("lambda:clean")
    dependsOn("cdk:clean")
}

tasks.register("deploy") {
    dependsOn("build")
    dependsOn("cdk:deploy")
}