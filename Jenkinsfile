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
        stage('Deploy') {
            steps {
                sh '''
            set -e

            docker create \
              --name go-hotfix-builder \
              golang:1.23

            docker cp main.go go-hotfix-builder:/src-main.go
            docker cp go.mod go-hotfix-builder:/go.mod

            docker exec go-hotfix-builder sh -c \
              'mkdir -p /src && mv /src-main.go /src/main.go && mv /go.mod /src/go.mod && \
               CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
               go build -ldflags="-s -w -X main.version=$VERSION" \
               -o /src/app /src/main.go'

            docker cp go-hotfix-builder:/src/app ./go_service_hotfix
            docker rm go-hotfix-builder

            docker cp go-service:/app ./go_service_backup

            docker cp go_service_hotfix go-service:/app

            docker restart go-service

            sleep 2

            curl -f http://host.docker.internal:8080
            '''
            }
        }
    }
}
