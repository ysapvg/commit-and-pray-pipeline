pipeline {
    agent any

    environment {
        IMAGE = 'go-service'
        DEPLOY_STARTED = 'false'
    }

    options {
        skipDefaultCheckout()
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
                          -t ${IMAGE}:${VERSION} .
                    """
                }
            }
        }

        // Replace the binary in the existing container
        stage('Deploy') {
            steps {
                script {
                    env.DEPLOY_STARTED = 'true'

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
                    '''
                }
            }
        }

        // Verify that the expected version is running
        stage('Verify') {
            steps {
                sh '''
                    set -e

                    sleep 2

                    RESPONSE=$(curl -fsS http://host.docker.internal:8080)

                    echo "Response: $RESPONSE"

                    # Simulate a deployment failure
                    exit 1

                    if [ "$RESPONSE" != "Hello, DevOps! version=$VERSION" ]; then
                        echo "Verification failed"
                        exit 1
                    fi
                '''
            }
        }
    }

    post {
        failure {
            script {
                if (env.DEPLOY_STARTED == 'true') {
                    echo 'Deployment failed. Rolling back...'

                    sh '''
                        docker cp go_service_backup go-service:/app
                        docker restart go-service

                        sleep 2

                        curl -fsS http://host.docker.internal:8080
                    '''
                }
            }
        }
    }
}
