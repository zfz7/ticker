import com.github.gradle.node.NodeExtension
import com.github.gradle.node.yarn.task.YarnTask

plugins {
    id("com.github.node-gradle.node")
}

configure<NodeExtension> {
    version.set(rootProject.extra["nodeVersion"] as String)
    yarnVersion.set(rootProject.extra["yarnVersion"] as String)
    download.set(true)
}

val install = tasks.register<YarnTask>("install") {
    inputs.file(file("$projectDir/yarn.lock"))
    inputs.file(file("$projectDir/package.json"))
    outputs.dir(file("$projectDir/node_modules"))
    args.set(listOf("install"))
}

tasks.register<YarnTask>("test") {
    environment.set(mapOf("CI" to "true"))
    dependsOn(install)
    args.set(listOf("test"))
}

tasks.register<YarnTask>("build") {
    dependsOn(install)
    dependsOn(":lambda:build")
    mustRunAfter("test")
    inputs.dir(file("$projectDir"))
    outputs.dir(file("$projectDir/build"))
    args.set(listOf("build"))
}

tasks.register<YarnTask>("deploy") {
    dependsOn(":build")
    args.set(listOf("deploy"))
}

task<Delete>("clean") {
    delete(project(":cdk").buildDir)
    delete(files("$projectDir/cdk.out"))
}