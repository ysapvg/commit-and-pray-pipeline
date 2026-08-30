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
    }
}