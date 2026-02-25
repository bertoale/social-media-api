pipeline {
    agent any

    environment {
        IMAGE_NAME = "albertxp/social-media-api"
        IMAGE_TAG  = "${env.BUILD_NUMBER}"
        TARGET_HOST = "192.168.56.119"
        TARGET_PATH = "/home/opt/social-media-api"
    }

    stages {

        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/bertoale/social-media-api.git'
            }
        }

        stage('Build Image') {
            steps {
                sh "docker build -t $IMAGE_NAME:$IMAGE_TAG ."
                sh "docker tag $IMAGE_NAME:$IMAGE_TAG $IMAGE_NAME:latest"
            }
        }

        stage('Login Docker Hub') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-login',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh '''
                    echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin
                    '''
                }
            }
        }

        stage('Push Image') {
            steps {
                sh "docker push $IMAGE_NAME:$IMAGE_TAG"
                sh "docker push $IMAGE_NAME:latest"
            }
        }

        stage('Deploy to Production') {
            steps {
                withCredentials([file(credentialsId: 'sosmed-env', variable: 'ENV_FILE')]) {

                    sshagent(['vm1-ssh-key']) {

                        sh """
                        scp -o StrictHostKeyChecking=no \
                            \$ENV_FILE albert@192.168.56.119:/opt/social-media-api/.env

                        ssh -o StrictHostKeyChecking=no \
                            albert@192.168.56.119 '
                            cd /home/opt/social-media-api &&
                            export IMAGE_TAG=${IMAGE_TAG} &&
                            docker pull ${IMAGE_NAME}:${IMAGE_TAG} &&
                            docker compose down &&
                            docker compose up -d
                            '
                        """
                    }
                }
            }
        }
    }

    post {
        success {
            echo "✅ Deployment Success - Version $IMAGE_TAG"
        }
        failure {
            echo "❌ Deployment Failed"
        }
    }
}