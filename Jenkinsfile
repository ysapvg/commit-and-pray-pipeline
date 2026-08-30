pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Test') {
            steps {
                sh '''
                    docker build \
                      --target builder \
                      -t go-service-test .

                    docker run --rm \
                      --entrypoint go \
                      go-service-test \
                      test ./...
                '''
            }
        }
        stage('Build Image') {
            steps {
                script {
                    env.VERSION = sh(
                script: 'git rev-parse --short HEAD',
                returnStdout: true
            ).trim()

                    sh """
                docker build \
                  --build-arg VERSION=${VERSION} \
                  -t go-service:${VERSION} .
            """
                }
            }
        }
    }
}
