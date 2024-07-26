pipeline {
    agent any
    parameters {
        string(name: 'PRJ_NAME', defaultValue: 'go-service-template', description: 'Project name')
    }
    stages {
        stage('build image') {
            steps {
                script {
                    withCredentials([sshUserPrivateKey(credentialsId: 'bitbucket', keyFileVariable: 'PRIVATE_KEY_FILE', usernameVariable: 'USERNAME')]) {
                        def imageName = params.PRJ_NAME.toLowerCase() + '-' + BRANCH_NAME.toLowerCase();
                        echo "${imageName}"
                        sh "cp ${PRIVATE_KEY_FILE} ./id_rsa"
                        sh "docker build --build-arg ssh_prv_key=./id_rsa --file docker/Dockerfile -t nexus.io.xunison.com/repository/${imageName}:latest ."
                        sh "docker tag nexus.io.xunison.com/repository/${imageName}:latest nexus.io.xunison.com/repository/${imageName}:${BUILD_NUMBER}-${GIT_COMMIT}"
                        sh "docker tag nexus.io.xunison.com/repository/${imageName}:latest harbor.us.io.xunison.com/xcpem/${imageName}:${BUILD_NUMBER}-${GIT_COMMIT}"
                        sh "docker tag nexus.io.xunison.com/repository/${imageName}:latest harbor.us.io.xunison.com/xcpem/${imageName}:latest"
                        echo env.BRANCH_NAME
                        sh "printenv"
                    }
                }
            }
        }
        stage('push image') {
            steps {
                script {
                    def imageName = params.PRJ_NAME.toLowerCase() + '-' + BRANCH_NAME.toLowerCase();
                    echo env.BRANCH_NAME
                    withCredentials([usernamePassword(credentialsId: 'f895fdab-0fbb-4fb9-a37b-84a5934ccc91', passwordVariable: 'password', usernameVariable: 'user')]) {
                        sh "docker login -u ${user} -p ${password} harbor.us.io.xunison.com" 
                        sh "docker push harbor.us.io.xunison.com/xcpem/${imageName}:latest"
                        sh "docker push harbor.us.io.xunison.com/xcpem/${imageName}:${BUILD_NUMBER}-${GIT_COMMIT}"
                    }
                    withCredentials([usernamePassword(credentialsId: 'jenkins-nexus', passwordVariable: 'password', usernameVariable: 'user')]) {
                        sh "docker login -u ${user} -p ${password} nexus.io.xunison.com" 
                        sh "docker push nexus.io.xunison.com/repository/${imageName}:latest"
                        sh "docker push nexus.io.xunison.com/repository/${imageName}:${BUILD_NUMBER}-${GIT_COMMIT}"
                    }
                }
            }
        }
        stage('build-push to ECR') {
            steps {
                script {
                    def scmVars = checkout scm
                    def awsRegistry = "709260118266.dkr.ecr.us-west-2.amazonaws.com"
                    withCredentials([sshUserPrivateKey(credentialsId: 'bitbucket', keyFileVariable: 'PRIVATE_KEY_FILE', usernameVariable: 'USERNAME')]) {
                        sh "cp ${PRIVATE_KEY_FILE} ./id_rsa"
                        sh "docker build --build-arg ssh_prv_key=id_rsa --file docker/Dockerfile -t ${awsRegistry}/iot:${env.BRANCH_NAME} -t ${awsRegistry}/iot:${scmVars.GIT_COMMIT} ."
                    }
                    docker.withRegistry("https://${awsRegistry}", "ecr:us-west-2:ecr-credentials") {
                        sh "docker push ${awsRegistry}/iot:${env.BRANCH_NAME}"
                        sh "docker push ${awsRegistry}/iot:${scmVars.GIT_COMMIT}"
                    }
                }
            }
        }
    }
}

