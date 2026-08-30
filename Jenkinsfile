pipeline {
    agent any

    stages {

        // Get the source code from GitHub
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        // Run the Go tests before building
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

        // Build and tag the Docker image using the Git commit hash
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

        // Replace the binary in the existing container and verify the deployment
        stage('Deploy') {
            steps {
                sh '''
                    set -e

                    # Get the binary from the newly built image
                    docker create --name hotfix-source ${IMAGE}:${VERSION}
                    docker cp hotfix-source:/app ./go_service_hotfix
                    docker rm hotfix-source

                    # Backup current binary
                    docker cp go-service:/app ./go_service_backup

                    # Replace binary and restart
                    docker cp go_service_hotfix go-service:/app
                    docker restart go-service

                    sleep 2
                    curl -f http://host.docker.internal:8080
                '''
            }
        }
    }
}