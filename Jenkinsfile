pipeline {
    agent any

    environment {
        IMAGE = 'go-service'
    }

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

        stage('Deploy') {
            steps {
                script {
                    try {
                        sh '''
                    set -e

                    # Get the binary from the new image
                    docker create --name hotfix-source ${IMAGE}:${VERSION}
                    docker cp hotfix-source:/app ./go_service_hotfix
                    docker rm hotfix-source

                    # Backup current binary
                    docker cp go-service:/app ./go_service_backup

                    # Replace binary and restart
                    docker cp go_service_hotfix go-service:/app
                    docker restart go-service

                    sleep 2
                    
                    # Simulate a failure for testing rollback
                    exit 1 

                    # Verify the new version
                    curl -f http://host.docker.internal:8080
                '''
            } catch (err) {
                        echo 'Deployment failed. Rolling back...'

                        sh '''
                    docker cp go_service_backup go-service:/app
                    docker restart go-service
                    sleep 2
                    curl -f http://host.docker.internal:8080
                '''

                        throw err
                    }
                }
            }
        }
    }
}
